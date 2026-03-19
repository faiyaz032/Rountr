package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const BACKEND_ADDR = "localhost:7777"

func main() {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from ", conn.RemoteAddr())

	backendConn, err := net.Dial("tcp", BACKEND_ADDR)
	if err != nil {
		log.Println("Error connecting to backend: ", err)
		return
	}
	defer backendConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(backendConn, conn)
		if err != nil {
			log.Println("Error copying to backend:", err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(conn, backendConn)
		if err != nil {
			log.Println("Error copying from backend:", err)
		}
		done <- struct{}{}
	}()

	<-done
	<-done
	fmt.Println("Connection from ", conn.RemoteAddr(), " closed")
}
