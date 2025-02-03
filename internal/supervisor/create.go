package supervisor

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ab-testing-service/internal/proxy"
)

type VirtualHostHandler struct {
	proxies     map[string]*proxy.Proxy // host -> proxy mapping
	pathProxies map[string]*proxy.Proxy // path -> proxy mapping
}

func (vh *VirtualHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// First try path-based routing
	path := r.URL.Path
	log.Printf("host: %s, path request: %s\n", r.Host, r.URL.Path)
	for pathKey, p := range vh.pathProxies {
		if strings.HasPrefix(path, "/"+pathKey) {
			// Remove the path key from the request path before proxying
			r.URL.Path = strings.TrimPrefix(path, "/"+pathKey)
			log.Printf("path: %s, key: %s, path request: %s\n", path, pathKey, r.URL.Path)
			p.ServeHTTP(w, r)
			return
		}
	}

	// If no path match, try host-based routing
	if r.Host == "" {
		http.Error(w, "host header is empty", http.StatusBadRequest)
		return
	}

	log.Printf("method: %s, host request: %s\n", r.Method, r.Host)
	host := strings.Split(r.Host, ":")[0]
	if p, ok := vh.proxies[host]; ok {
		p.ServeHTTP(w, r)
	} else {
		http.Error(w, fmt.Sprintf("host %s not found", host), http.StatusNotFound)
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

	// Initialize virtual handler if needed
	if s.virtualHandler == nil {
		s.virtualHandler = &VirtualHostHandler{
			proxies:     make(map[string]*proxy.Proxy),
			pathProxies: make(map[string]*proxy.Proxy),
		}
	}

	p, err := proxy.NewProxy(cfg)
	if err != nil {
		return err
	}

	instance := &ProxyInstance{
		Proxy: p,
	}
	s.proxies[cfg.ID] = instance

	// Handle path-based routing
	if cfg.PathKey != "" {
		if _, exists := s.virtualHandler.pathProxies[cfg.PathKey]; exists {
			return fmt.Errorf("proxy with path key %s already exists", cfg.PathKey)
		}
	}

	// Handle host-based routing
	if _, exists := s.virtualHandler.proxies[host]; exists {
		return fmt.Errorf("proxy with host %s already exists", host)
	}

	if cfg.PathKey != "" {
		log.Printf("path key: %s\n", cfg.PathKey)
		s.virtualHandler.pathProxies[cfg.PathKey] = p
	} else {
		log.Printf("host: %s\n", host)
		s.virtualHandler.proxies[host] = p
	}

	instance.Started = true

	// Start the server if not already running
	if s.server == nil {
		s.server = &http.Server{
			Addr:    ":80", // Default port, should be configurable
			Handler: s.virtualHandler,
		}

		go func() {
			if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// Log error but don't stop the server
				log.Printf("main HTTP server error: %v\n", err)
			}
		}()
	}

	return nil
}
