package proxy

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/faiyaz032/rountr/balancer"
)

func HandleConnection(conn net.Conn, lb balancer.LoadBalancer) {
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
