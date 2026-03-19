package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type LoadBalancer struct {
	servers []string
	current int
	mu      sync.Mutex
}

func main() {
	loadBalancer := NewLoadBalancer([]string{
		"localhost:7777",
		"localhost:7778",
	})

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("Rountr listening on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn, loadBalancer)
	}
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
	}
}

func (lb *LoadBalancer) GetNextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)
	return server
}

func handleConnection(conn net.Conn, lb *LoadBalancer) {
	defer conn.Close()
	fmt.Println("New connection from ", conn.RemoteAddr())

	serverAddr := lb.GetNextServer()
	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println("Error connecting to backend: ", err)
		return
	}
	defer serverConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(serverConn, conn)
		if err != nil {
			log.Println("Error copying to backend:", err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(conn, serverConn)
		if err != nil {
			log.Println("Error copying from backend:", err)
		}
		done <- struct{}{}
	}()

	<-done
	<-done
	fmt.Println("Connection from ", conn.RemoteAddr(), " closed")
}
