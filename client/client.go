package main

import (
	"context"
	"log"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Rbd3178/grpcMessageBoard/messageBoard"
)

func main() {
	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessageBoardClient(conn)

	message := pb.Message{
		Author: "IlyaV",
		Title:  "The very first",
		Body:   "Test",
	}

	response, err := client.PostMessage(context.Background(), &message)
	if err != nil {
		log.Fatalf("Failed to call PostMessage: %v", err)
	}

	log.Printf("Response from the server. Id: %d Author: %s Title: %s Body: %s", response.Id, response.Author, response.Body, response.Title)
}
