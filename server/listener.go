package server

import (
	"log"
	"net"

	"github.com/faiyaz032/rountr/balancer"
	"github.com/faiyaz032/rountr/proxy"
)

func Start(address string, lb balancer.LoadBalancer) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("Load balancer listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go proxy.HandleConnection(conn, lb)
	}
}
