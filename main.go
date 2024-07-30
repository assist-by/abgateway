package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Lux-N-Sal/autro-signal/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSignalServiceServer
}

func (s *server) SendSignal(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	log.Printf("Received signal: %v", req)
	fmt.Printf("Signal: %s\n", req.Signal)
	fmt.Printf("Timestamp: %d\n", req.Timestamp)
	fmt.Printf("Price: %s\n", req.Price)
	fmt.Printf("Long conditions: %v\n", req.Conditions.Long)
	fmt.Printf("Short conditions: %v\n", req.Conditions.Short)

	return &pb.SignalResponse{
		Success: true,
		Message: "Signal received and processed",
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSignalServiceServer(s, &server{})

	log.Printf("API Gateway gRPC server listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
