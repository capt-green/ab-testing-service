package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ab-testing-service/internal/models"
	"github.com/ab-testing-service/internal/proxy"
	"github.com/ab-testing-service/internal/storage"
)

type UpdateTargetsRequest struct {
	Targets []struct {
		URL      string  `json:"url" binding:"required"`
		Weight   float64 `json:"weight" binding:"required,min=0,max=1"`
		IsActive bool    `json:"is_active"`
	} `json:"targets"`
	Condition *RouteCondition `json:"condition,omitempty"`
}

func (s *Server) updateProxyTargets(c *gin.Context) {
	proxyID := c.Param("id")

	req, err := s.parseAndValidateRequest(c)
	if err != nil {
		return // Error already sent to client
	}

	currentProxy, err := s.getCurrentProxy(c, proxyID)
	if err != nil {
		return // Error already sent to client
	}

	targets := s.convertToTargetModels(proxyID, req)
	condition := s.convertToConditionModels(targets, req)

	if err := s.executeTransaction(c, proxyID, currentProxy, targets, condition); err != nil {
		return // Error already sent to client
	}

	if err := s.updateSupervisor(c, proxyID, currentProxy, targets, condition); err != nil {
		return // Error already sent to client
	}

	c.JSON(http.StatusOK, gin.H{"message": "proxy updated successfully"})
}

// Request parsing and validation
func (s *Server) parseAndValidateRequest(c *gin.Context) (UpdateTargetsRequest, error) {
	var req UpdateTargetsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return req, err
	}

	if err := s.validateCondition(c, &req); err != nil {
		return req, err
	}

	return req, nil
}

func (s *Server) validateCondition(c *gin.Context, req *UpdateTargetsRequest) error {
	if req.Condition == nil {
		return nil
	}

	if err := s.validateConditionFields(req.Condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	//if err := s.validateConditionTargets(req); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return err
	//}

	return nil
}

func (s *Server) validateConditionFields(condition *RouteCondition) error {
	if !models.ConditionType(condition.Type).IsValid() {
		return errors.New("invalid condition type")
	}
	if condition.ParamName == "" {
		return errors.New("param_name is required for query_param condition")
	}
	if len(condition.Values) == 0 {
		return errors.New("values map is required for query_param condition")
	}
	return nil
}

func (s *Server) validateConditionTargets(req *UpdateTargetsRequest) error {
	targetIDs := make(map[string]bool)
	for _, target := range req.Targets {
		targetIDs[target.URL] = true // Assuming URL is unique we have no id of new target here
	}

	for _, targetID := range req.Condition.Values {
		if !targetIDs[targetID] {
			return fmt.Errorf("target ID %s in condition not found in targets", targetID)
		}
	}

	if req.Condition.Default != "" && !targetIDs[req.Condition.Default] {
		return errors.New("default target ID not found in targets")
	}

	return nil
}

// Helper functions
func (s *Server) getCurrentProxy(c *gin.Context, proxyID string) (*models.Proxy, error) {
	p, err := s.storage.GetProxy(c.Request.Context(), proxyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to get current proxy state: %v", err)})
		return nil, err
	}
	return p, nil
}

func (s *Server) convertToTargetModels(proxyID string, req UpdateTargetsRequest) []models.Target {
	targets := make([]models.Target, len(req.Targets))
	for i, t := range req.Targets {
		targets[i] = models.Target{
			ID:       uuid.New().String(),
			ProxyID:  proxyID,
			URL:      t.URL,
			Weight:   t.Weight,
			IsActive: t.IsActive,
		}
	}
	return targets
}

func (s *Server) convertToConditionModels(targets []models.Target, req UpdateTargetsRequest) *models.RouteCondition {
	if req.Condition == nil || req.Condition.Type == "" {
		return nil
	}

	conditionValues := make(map[string]string, len(req.Condition.Values))
	for i, v := range req.Condition.Values {
		conditionValues[targets[i].ID] = v
	}
	return &models.RouteCondition{
		Type:      models.ConditionType(req.Condition.Type),
		ParamName: req.Condition.ParamName,
		Values:    conditionValues,
		Default:   req.Condition.Default,
	}
}

func (s *Server) getUserID(c *gin.Context) *string {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*models.User); ok {
			return &u.ID
		}
	}
	return nil
}

// Transaction handling
func (s *Server) executeTransaction(c *gin.Context, proxyID string, currentProxy *models.Proxy,
	targets []models.Target, condition *models.RouteCondition) error {

	tx, err := s.storage.BeginTx(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to start transaction: %v", err)})
		return err
	}
	defer tx.Rollback()

	userID := s.getUserID(c)

	if err := s.recordChanges(c, tx, proxyID, currentProxy, targets, condition, userID); err != nil {
		return err
	}

	if err := s.updateStorage(c, tx, proxyID, targets, condition); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to commit transaction: %v", err)})
		return err
	}

	return nil
}

func (s *Server) recordChanges(c *gin.Context, tx *storage.Tx, proxyID string,
	currentProxy *models.Proxy, targets []models.Target,
	condition *models.RouteCondition, userID *string) error {

	if err := s.storage.RecordProxyChange(
		c.Request.Context(),
		tx,
		proxyID,
		models.ChangeTypeTargetsUpdate,
		currentProxy.Targets,
		targets,
		userID,
	); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to record target changes: %v", err)})
		return err
	}

	if condition != nil {
		if err := s.storage.RecordProxyChange(
			c.Request.Context(),
			tx,
			proxyID,
			models.ChangeTypeConditionUpdate,
			currentProxy.Condition,
			condition,
			userID,
		); err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("failed to record condition changes: %v", err)})
			return err
		}
	}

	return nil
}

func (s *Server) updateStorage(c *gin.Context, tx *storage.Tx, proxyID string,
	targets []models.Target, condition *models.RouteCondition) error {

	if err := s.storage.UpdateTargetsWithTx(c.Request.Context(), tx, proxyID, targets); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to update targets: %v", err)})
		return err
	}

	if condition != nil {
		if err := s.storage.UpdateProxyConditionWithTx(c.Request.Context(), tx, proxyID, condition); err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("failed to update condition: %v", err)})
			return err
		}
	}

	return nil
}

// Supervisor update
func (s *Server) updateSupervisor(c *gin.Context, proxyID string, currentProxy *models.Proxy,
	targets []models.Target, condition *models.RouteCondition) error {

	config := s.buildProxyConfig(proxyID, currentProxy, targets, condition)

	if err := s.supervisor.UpdateProxyTargets(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to update proxy targets: %v", err)})
		return err
	}

	return nil
}

func (s *Server) buildProxyConfig(proxyID string, currentProxy *models.Proxy,
	targets []models.Target, condition *models.RouteCondition) proxy.Config {

	config := proxy.Config{
		ID:        proxyID,
		ListenURL: currentProxy.ListenURL,
		Mode:      models.ProxyMode(currentProxy.Mode),
		Targets:   s.convertToConfigTargets(targets),
	}

	if condition != nil {
		config.Condition = &proxy.Condition{
			Type:      condition.Type,
			ParamName: condition.ParamName,
			Values:    condition.Values,
			Default:   condition.Default,
		}
	}

	return config
}

func (s *Server) convertToConfigTargets(targets []models.Target) []proxy.Target {
	configTargets := make([]proxy.Target, len(targets))
	for i, t := range targets {
		configTargets[i] = proxy.Target{
			ID:       t.ID,
			URL:      t.URL,
			Weight:   t.Weight,
			IsActive: t.IsActive,
		}
	}
	return configTargets
}
