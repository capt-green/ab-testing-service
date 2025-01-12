package supervisor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ab-testing-service/internal/proxy"
)

type VirtualHostHandler struct {
	proxies map[string]*proxy.Proxy
}

func (vh *VirtualHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	if p, ok := vh.proxies[host]; ok {
		p.ServeHTTP(w, r)
	} else {
		http.Error(w, "Host not found", http.StatusNotFound)
	}
}

func (s *Supervisor) CreateProxy(cfg proxy.Config) error {
	// Extract hostname from ListenURL
	host := strings.Split(cfg.ListenURL, ":")[0]
	if host == "" {
		return fmt.Errorf("invalid listen URL: %s", cfg.ListenURL)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if proxy with same host already exists
	if s.virtualHandler != nil {
		if _, exists := s.virtualHandler.proxies[host]; exists {
			return fmt.Errorf("proxy with host %s already exists", host)
		}
	}

	p, err := proxy.NewProxy(cfg)
	if err != nil {
		return err
	}

	instance := &ProxyInstance{
		Proxy: p,
	}

	// Initialize virtual host handler if not exists
	if s.virtualHandler == nil {
		s.virtualHandler = &VirtualHostHandler{
			proxies: make(map[string]*proxy.Proxy),
		}
	}

	// Add proxy to virtual host handler
	s.virtualHandler.proxies[host] = p
	s.proxies[cfg.ID] = instance

	// Start the server if not already running
	if s.server == nil {
		s.server = &http.Server{
			Addr:    ":80", // Default port, should be configurable
			Handler: s.virtualHandler,
		}

		go func() {
			if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// Log error but don't stop the server
				fmt.Printf("HTTP server error: %v\n", err)
			}
		}()
	}

	return nil
}
