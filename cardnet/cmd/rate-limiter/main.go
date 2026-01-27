package main

import (
	"log"
	"net"
	"os"

	pb "cardnet/internal/api/grpc"
	"cardnet/internal/rate"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Rate Limiter...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "60051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	limiter := rate.NewLimiter()
	handler := rate.NewHandler(limiter)

	grpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("Rate Limiter listening on :60051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("rate-limiter server failed: %v", err)
	}
}
