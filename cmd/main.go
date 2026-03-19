package main

import (
	"log"

	"github.com/faiyaz032/rountr/balancer"
	"github.com/faiyaz032/rountr/server"
)

func main() {
	lb := balancer.NewRoundRobin([]string{
		"localhost:7777",
		"localhost:7778",
	})

	if err := server.Start(":8888", lb); err != nil {
		log.Fatal(err)
	}
}
