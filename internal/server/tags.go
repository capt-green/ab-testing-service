package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UpdateTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

func (s *Server) updateProxyTags(c *gin.Context) {
	proxyID := c.Param("id")

	var req UpdateTagsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateProxyTags(c.Request.Context(), proxyID, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "proxy tags updated successfully"})
}

func (s *Server) getAllTags(c *gin.Context) {
	tags, err := s.storage.GetAllTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func (s *Server) getProxiesByTags(c *gin.Context) {
	var tags []string
	if tagsParam := c.Query("tags"); tagsParam != "" {
		tags = strings.Split(tagsParam, ",")
	}

	proxies, err := s.storage.GetProxiesByTags(c.Request.Context(), tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"proxies": proxies})
}
