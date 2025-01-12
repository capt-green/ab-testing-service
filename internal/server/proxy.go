package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ab-testing-service/internal/models"
	"github.com/ab-testing-service/internal/proxy"
)

type RouteCondition struct {
	Type      string   `json:"type" db:"type"`        // Type of condition: "header", "query", "cookie", "user_agent", "language"
	ParamName string   `json:"param_name" db:"param"` // Name of the parameter to check (for header, query, cookie)
	Values    []string `json:"values" db:"values"`    // List of parameter values to match targets
	Default   string   `json:"default" db:"default"`  // Default target ID if no match is found
}

type CreateProxyRequest struct {
	ListenURL string             `json:"listen_url" binding:"required"`
	Mode      string             `json:"mode" binding:"required"`
	Tags      []string           `json:"tags"`
	Targets   []CreateTargetSpec `json:"targets"`
	Condition *RouteCondition    `json:"condition,omitempty"`
}

type CreateTargetSpec struct {
	URL      string  `json:"url" binding:"required"`
	Weight   float64 `json:"weight" binding:"required,min=0,max=1"`
	IsActive bool    `json:"is_active"`
}

func (s *Server) createProxy(c *gin.Context) {
	var req CreateProxyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default mode if not specified
	if req.Mode == "" {
		req.Mode = string(models.ProxyModeReverse)
	}

	// Validate mode
	if req.Mode != string(models.ProxyModeReverse) && req.Mode != string(models.ProxyModeRedirect) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proxy mode"})
		return
	}

	// Create proxy model
	p := &models.Proxy{
		ListenURL: req.ListenURL,
		Mode:      req.Mode,
		Tags:      req.Tags,
	}

	// Convert targets
	if len(req.Targets) > 0 {
		p.Targets = make([]models.Target, len(req.Targets))
		for i, t := range req.Targets {
			p.Targets[i] = models.Target{
				ID:       uuid.New().String(),
				URL:      t.URL,
				Weight:   t.Weight,
				IsActive: t.IsActive,
			}
		}
	}

	// Convert condition
	if req.Condition.Type != "" {
		conditionType := models.ConditionType(req.Condition.Type)
		if !conditionType.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid condition type"})
			return
		}

		conditionValues := make(map[string]string, len(req.Condition.Values))
		for i, v := range req.Condition.Values {
			conditionValues[p.Targets[i].ID] = v
		}

		p.Condition = &models.RouteCondition{
			Type:      conditionType,
			ParamName: req.Condition.ParamName,
			Values:    conditionValues,
			Default:   req.Condition.Default,
		}
	}

	// Create proxy in storage -> postgres
	if err := s.storage.CreateProxy(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create proxy in storage",
			"details": err.Error(),
		})
		return
	}

	// Create proxy configuration for supervisor
	cfg := proxy.Config{
		ID:        p.ID,
		ListenURL: p.ListenURL,
		Mode:      models.ProxyMode(p.Mode),
	}

	// Convert targets to config format
	if len(p.Targets) > 0 {
		cfg.Targets = make([]proxy.Target, len(p.Targets))
		for i, t := range p.Targets {
			cfg.Targets[i] = proxy.Target{
				ID:       t.ID,
				URL:      t.URL,
				Weight:   t.Weight,
				IsActive: t.IsActive,
			}
		}
	}

	// Add condition if provided
	if p.Condition != nil {
		cfg.Condition = &proxy.Condition{
			Type:      p.Condition.Type,
			ParamName: p.Condition.ParamName,
			Values:    p.Condition.Values,
			Default:   p.Condition.Default,
		}
	}

	// Create proxy in supervisor -> start proxy server
	if err := s.supervisor.CreateProxy(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create proxy in supervisor",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, p)
}
