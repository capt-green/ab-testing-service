package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ab-testing-service/internal/config"
	"github.com/ab-testing-service/internal/storage"
	"github.com/ab-testing-service/internal/supervisor"
)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	supervisor *supervisor.Supervisor
	storage    *storage.Storage
	srv        *http.Server
}

func NewServer(cfg *config.Config, sup *supervisor.Supervisor, storage *storage.Storage) *Server {
	s := &Server{
		config:     cfg,
		supervisor: sup,
		storage:    storage,
	}

	s.setupRouter()
	return s
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

type GetProxyListRequest struct {
	Limit    int    `form:"limit,default=10"`
	Offset   int    `form:"offset,default=0"`
	SortBy   string `form:"sortBy"`
	SortDesc bool   `form:"sortDesc,default=false"`
}

func (s *Server) listProxies(c *gin.Context) {
	var req GetProxyListRequest
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

	proxies := s.supervisor.ListProxies(req.SortBy, req.SortDesc)
	
	// Calculate total count for pagination
	total := len(proxies)
	
	// Apply pagination
	end := req.Offset + req.Limit
	if end > total {
		end = total
	}
	if req.Offset > total {
		req.Offset = total
	}
	
	paginatedProxies := proxies[req.Offset:end]
	
	c.JSON(http.StatusOK, gin.H{
		"items": paginatedProxies,
		"total": total,
	})
}

func (s *Server) getProxy(c *gin.Context) {
	id := c.Param("id")
	proxy := s.supervisor.GetProxy(id)
	if proxy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "proxy not found"})
		return
	}

	c.JSON(http.StatusOK, proxy)
}

func (s *Server) deleteProxy(c *gin.Context) {
	id := c.Param("id")
	if err := s.supervisor.DeleteProxy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
