package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb_notification "github.com/Lux-N-Sal/autro-notification/proto"
	pb_signal "github.com/Lux-N-Sal/autro-signal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb_signal.UnimplementedSignalServiceServer
	notificationClient pb_notification.NotificationServiceClient
}

func (s *server) SendSignal(ctx context.Context, req *pb_signal.SignalRequest) (*pb_signal.SignalResponse, error) {
	log.Printf("Received signal: %v", req)
	description := fmt.Sprintf("Signal: %s for BTCUSDT at %v\n\n"+
		"LONG  - CASE 1: %t, CASE 2: %t, CASE 3: %t\n"+
		"SHORT - CASE 1: %t, CASE 2: %t, CASE 3: %t",
		req.Signal, time.Unix(req.Timestamp, 0),
		req.Conditions.Long[0], req.Conditions.Long[1], req.Conditions.Long[2],
		req.Conditions.Short[0], req.Conditions.Short[1], req.Conditions.Short[2])

	notificationReq := &pb_notification.NotificationRequest{
		Title:       fmt.Sprintf("New Signal: %s", req.Signal),
		Description: description,
		Signal:      req.Signal,
	}
	notificationResp, err := s.notificationClient.SendNotification(ctx, notificationReq)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return &pb_signal.SignalResponse{Success: false, Message: "Failed to send notification"}, nil
	}
	return &pb_signal.SignalResponse{
		Success: notificationResp.Success,
		Message: notificationResp.Message,
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	notificationServiceAddr := os.Getenv("NOTIFICATION_SERVICE_ADDR")
	if notificationServiceAddr == "" {
		notificationServiceAddr = "notification-service:50052"
	}

	notificationConn, err := grpc.NewClient(notificationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to notification service: %v", err)
	}
	defer notificationConn.Close()

	notificationClient := pb_notification.NewNotificationServiceClient(notificationConn)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb_signal.RegisterSignalServiceServer(s, &server{notificationClient: notificationClient})
	reflection.Register(s)

	log.Printf("API Gateway gRPC server listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
