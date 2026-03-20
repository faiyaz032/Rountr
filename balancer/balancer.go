package balancer

const (
	RoundRobinAlgo      = "round_robin"
	LeastConnectionAlgo = "least_connection"
)

func NewLoadBalancer(algo string, servers []string) LoadBalancer {
	switch algo {
	case "round_robin":
		return NewRoundRobin(servers)
	case "least_connection":
		return NewLeastConnections(servers)
	default:
		panic("unsupported load balancing algorithm: " + algo)
	}
}
