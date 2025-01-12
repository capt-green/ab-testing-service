package supervisor

import (
	"sort"

	"github.com/ab-testing-service/internal/proxy"
)

func (s *Supervisor) GetProxy(id string) *proxy.Proxy {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.proxies[id].Proxy
}

func (s *Supervisor) ListProxies(sortBy string, sortDesc bool) []proxy.Config {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var configs []proxy.Config
	for id, p := range s.proxies {
		tags := s.storage.GetTags(id)
		configs = append(configs, proxy.Config{
			ID:        id,
			ListenURL: p.Proxy.ListenURL,
			Mode:      p.Proxy.Mode,
			Targets:   p.Proxy.Targets,
			Condition: p.Proxy.Config.Condition,
			Tags:      tags,
		})
	}

	// Sort the configs based on the sortBy parameter
	if sortBy != "" {
		sort.Slice(configs, func(i, j int) bool {
			var result bool
			switch sortBy {
			case "id":
				result = configs[i].ID < configs[j].ID
			case "listen_url":
				result = configs[i].ListenURL < configs[j].ListenURL
			case "mode":
				result = configs[i].Mode < configs[j].Mode
			case "targets":
				result = len(configs[i].Targets) < len(configs[j].Targets)
			default:
				return !sortDesc // Default sort by ID
			}
			if sortDesc {
				return !result
			}
			return result
		})
	}

	return configs
}
