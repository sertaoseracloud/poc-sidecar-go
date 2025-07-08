package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	pb "sidecar/sidecar/proto"
)

const (
	authQueue = "auth_requests"
)

func sendAuthRequest(provider string) (string, error) {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://rabbitmq"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return "", err
	}
	ch, err := conn.Channel()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return "", err
	}
	corrID := uuid.New().String()
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return "", err
	}

	if err := ch.Publish("", authQueue, false, false,
		amqp.Publishing{
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
			Body:          []byte(provider),
		}); err != nil {
		return "", err
	}
	for msg := range msgs {
		if msg.CorrelationId == corrID {
			msg.Ack(false)
			return string(msg.Body), nil
		}
	}
	return "", nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Provider string `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := sendAuthRequest(req.Provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

// gRPC server

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) RequestAuth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	res, err := sendAuthRequest(in.Provider)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Json: res}, nil
}

func main() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "3000"
	}

	http.HandleFunc("/auth", httpHandler)
	go func() {
		log.Printf("HTTP server on %s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Fatalf("http failed: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	log.Println("gRPC server on 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc failed: %v", err)
	}
}
