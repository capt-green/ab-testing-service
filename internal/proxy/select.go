package proxy

import (
	"fmt"
	"math/rand"
	"net/http"
)

func (p *Proxy) selectTarget(r *http.Request) (*Target, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// First, try to get target from cookie
	if target := p.getTargetFromCookie(r); target != nil {
		return target, nil
	}

	// Then, check routing conditions if present
	if p.Config.Condition != nil {
		if target := p.getTargetByCondition(r); target != nil {
			return target, nil
		}
		return nil, fmt.Errorf("no matching target found")
	}

	// Fall back to weighted random selection if no condition matches or no condition is set
	var totalWeight float64
	var activeTargets []Target

	for _, target := range p.Targets {
		if target.IsActive {
			totalWeight += target.Weight
			activeTargets = append(activeTargets, target)
		}
	}

	if len(activeTargets) == 0 {
		return nil, fmt.Errorf("no active targets available")
	}

	// Choose a random target based on weights
	rnd := rand.Float64() * totalWeight
	var cumulativeWeight float64

	for _, target := range activeTargets {
		cumulativeWeight += target.Weight
		if rnd <= cumulativeWeight {
			return &target, nil
		}
	}

	// Fallback to the first active target if something goes wrong with the random selection
	return &activeTargets[0], nil
}
