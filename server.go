package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalln("Failed to listen on port 8090: ", err)
	}

	server := grpc.NewServer()

	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serve on port 8090: ", err)
	}
}