package balancer

type LoadBalancer interface {
	GetNextServer() string
}
