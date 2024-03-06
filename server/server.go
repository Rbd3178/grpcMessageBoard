package main

import (
	"context"
	"log"
	"net"
	"sync"
	"fmt"

	pb "github.com/Rbd3178/grpcMessageBoard/messageBoard"
	"google.golang.org/grpc"
)

type messageBoardServer struct {
	pb.UnimplementedMessageBoardServer
	messages []*pb.Message
	maxSize  int
	mu       sync.RWMutex
}

func newServer(maxSize int) *messageBoardServer {
	return &messageBoardServer{messages: make([]*pb.Message, 0), maxSize: maxSize}
}

func (s *messageBoardServer) PostMessage(c context.Context, m *pb.Message) (*pb.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.messages) == 0 {
		m.Id = 1
	} else {
		m.Id = s.messages[len(s.messages)-1].Id + 1
	}

	if len(s.messages) >= s.maxSize {
		s.messages = s.messages[1:]
	}
	s.messages = append(s.messages, m)
	return m, nil
}

func (s *messageBoardServer) GetLatestMessages(r *pb.GetLatestRequest, stream pb.MessageBoard_GetLatestMessagesServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := len(s.messages)
	start := count - int(r.Amount)
	if start < 0 {
		start = 0
	}

	for _, message := range s.messages[start:] {
		err := stream.Send(message)
		if err != nil {
			return err
		}
	}

	return nil
}

const port = 8090

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageBoardServer(grpcServer, newServer(10))

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve on port %d: %v", port, err)
	}
}
