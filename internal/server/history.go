package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetProxyChangesRequest struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

func (s *Server) getProxyChanges(c *gin.Context) {
	proxyID := c.Param("id")

	var req GetProxyChangesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and cap the limit
	if req.Limit <= 0 {
		req.Limit = 10
	} else if req.Limit > 100 {
		req.Limit = 100
	}

	// Ensure offset is not negative
	if req.Offset < 0 {
		req.Offset = 0
	}

	changes, err := s.storage.GetProxyChanges(c.Request.Context(), proxyID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"changes": changes,
		"pagination": gin.H{
			"limit":  req.Limit,
			"offset": req.Offset,
			"total":  len(changes),
		},
	})
}
