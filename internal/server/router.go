package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ab-testing-service/internal/middleware"
)

func (s *Server) setupRouter() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", s.login)
		auth.POST("/register", s.register)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(s.config))
	{
		api.GET("/proxies", s.listProxies)
		api.POST("/proxies", s.createProxy)
		api.GET("/proxies/:id", s.getProxy)
		api.DELETE("/proxies/:id", s.deleteProxy)
		api.PUT("/proxies/:id/targets", s.updateProxyTargets)
		api.GET("/proxies/:id/changes", s.getProxyChanges)

		// Tag management
		api.GET("/tags", s.getAllTags)
		api.GET("/proxies/by-tags", s.getProxiesByTags)
		api.PUT("/proxies/:id/tags", s.updateProxyTags)

		// Metrics
		api.GET("/metrics", func(c *gin.Context) {
			handler := promhttp.Handler()
			handler.ServeHTTP(c.Writer, c.Request)
		})
	}

	s.router = r
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: r,
	}
}
