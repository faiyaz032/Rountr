package balancer

import "sync"

type LeastConnections struct {
	servers []string
	counts  map[string]int
	mu      sync.Mutex
}

func NewLeastConnections(servers []string) *LeastConnections {
	counts := make(map[string]int)
	for _, s := range servers {
		counts[s] = 0
	}
	return &LeastConnections{
		servers: servers,
		counts:  counts,
	}
}

func (lc *LeastConnections) GetNextServer() string {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var selected string
	min := int(^uint(0) >> 1)

	for _, s := range lc.servers {
		if lc.counts[s] < min {
			min = lc.counts[s]
			selected = s
		}
	}

	lc.counts[selected]++
	return selected
}
