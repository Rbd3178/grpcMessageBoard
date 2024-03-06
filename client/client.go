package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Rbd3178/grpcMessageBoard/messageBoard"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	address, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	address = strings.TrimSuffix(address, "\n")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewMessageBoardClient(conn)
	fmt.Printf("Successfully connected to message board at %s. Type 'help' for available commands.\n", address)

	name := "anonymous"

	for {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		// Execute commands
		switch parts[0] {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  help - Show available commands")
			fmt.Println("  showname - Print current name. Default is \"anonymous\"")
			fmt.Println("  setname [name] - Set your name")
			fmt.Println("  postmsg - Start the process of posting a message")
			fmt.Println("  showmsg [count] - Print <count> latest messages")
			fmt.Println("  exit - Exit the application")
			fmt.Println()

		case "showname":
			fmt.Printf("Current name is \"%s\"\n\n", name)

		case "setname":
			if len(parts) < 2 {
				fmt.Printf("Not enough arguments\n\n")
				continue
			}
			name = parts[1]
			fmt.Printf("Set name to \"%s\"\n\n", parts[1])

		case "postmsg":
			fmt.Print("Enter title: ")
			title, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			fmt.Print("Enter message: ")
			body, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			message := pb.Message{
				Author: name,
				Title:  strings.TrimSuffix(title, "\n"),
				Body:   strings.TrimSuffix(body, "\n"),
			}
			response, err := client.PostMessage(context.Background(), &message)
			if err != nil {
				fmt.Printf("Failed to post the message: %v\n\n", err)
			}
			fmt.Println("Message posted successfully:")
			fmt.Printf("Message №%d by %s: %s\n %s\n\n", response.Id, response.Author, response.Title, response.Body)

		case "showmsg":
			if len(parts) < 2 {
				fmt.Printf("Not enough arguments\n\n")
				continue
			}

			count, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Incorrect argument format: %s\n\n", err)
				continue
			}
			request := pb.GetLatestRequest{
				Amount: int32(count),
			}
			stream, err := client.GetLatestMessages(context.Background(), &request)
			if err != nil {
				fmt.Printf("Failed to get messages: %v\n\n", err)
				continue
			}

			for {
				message, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Error when getting the messages: %v\n\n", err)
					break
				}
				fmt.Printf("Message №%d by %s: %s\n %s\n\n", message.Id, message.Author, message.Title, message.Body)
			}

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Printf("Unknown command. Type 'help' for available commands.\n\n")
		}
	}
}
