package balancer

import "sync"

type RoundRobin struct {
	servers []string
	current int
	mu      sync.Mutex
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{
		servers: servers,
	}
}

func (rr *RoundRobin) GetNextServer() string {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	server := rr.servers[rr.current]
	rr.current = (rr.current + 1) % len(rr.servers)
	return server
}
