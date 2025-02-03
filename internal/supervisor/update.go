package supervisor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ab-testing-service/internal/proxy"
)

func (s *Supervisor) DeleteProxy(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get the instance before deleting
	instance, exists := s.proxies[id]
	if exists && instance.Proxy != nil {
		// Remove from virtual host handler
		host := strings.Split(instance.Proxy.Config.ListenURL, ":")[0]
		if s.virtualHandler != nil {
			delete(s.virtualHandler.proxies, host)
		}
	}

	// Remove from proxies map
	delete(s.proxies, id)

	return s.storage.InvalidateProxyCache(ctx, id)
}

// handleProxyUpdate is called when a proxy settings change notification is received
func (s *Supervisor) handleProxyUpdate(ctx context.Context, proxyID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get the latest config from storage
	cfg, err := s.storage.GetProxyConfig(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to get proxy config: %w", err)
	}

	// Update the proxy
	return s.UpdateProxyTargets(ctx, cfg)
}

func (s *Supervisor) UpdateProxyTargets(ctx context.Context, cfg proxy.Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	instance, exists := s.proxies[cfg.ID]
	if !exists {
		return fmt.Errorf("proxy %s not found", cfg.ID)
	}

	// Get old and new hosts
	oldHost := ""
	if instance.Proxy != nil {
		oldHost = strings.Split(instance.Proxy.Config.ListenURL, ":")[0]
	}
	newHost := strings.Split(cfg.ListenURL, ":")[0]

	// Create new proxy with updated config
	newProxy, err := proxy.NewProxy(cfg)
	if err != nil {
		return fmt.Errorf("failed to create new proxy: %w", err)
	}

	// Update virtual host handler
	if s.virtualHandler != nil {
		// Remove old host if it changed
		if oldHost != "" && oldHost != newHost {
			delete(s.virtualHandler.proxies, oldHost)
		}
		// Add new host
		s.virtualHandler.proxies[newHost] = newProxy
	}

	// Update the instance
	s.proxies[cfg.ID] = &ProxyInstance{
		Proxy:   newProxy,
		Started: true,
	}

	if err := s.storage.InvalidateProxyCache(ctx, cfg.ID); err != nil {
		return fmt.Errorf("failed to invalidate proxy cache: %w", err)
	}

	if err := s.storage.SaveProxyConfig(ctx, cfg); err != nil {
		return fmt.Errorf("failed to update proxy config: %w", err)
	}

	// Publish change to other instances
	if err := s.pubsub.PublishSettingsChange(ctx, cfg.ID); err != nil {
		log.Printf("Failed to publish settings change: %v", err)
	}

	return nil
}
