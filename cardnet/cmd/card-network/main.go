package main

import (
	"cardnet/internal/issuer/config"
	"cardnet/internal/issuer/routing"
	"log"
	"net"
	"os"

	issuerpb "cardnet/internal/api/issuer"
	"cardnet/internal/network"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const issuerConfigPath = "internal/issuer/config/issuers.yaml"

func main() {
	// Load BIN routing table (required for ResolveIssuer to work)
	table, err := config.LoadIssuers(issuerConfigPath)
	if err != nil {
		log.Fatalf("failed to load issuer config: %v", err)
	}
	routing.Load(table)
	log.Println("Loaded BIN routing table")

	port := os.Getenv("PORT")
	if port == "" {
		port = "60053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	issuerpb.RegisterIssuerServiceServer(server, network.NewHandler())
	reflection.Register(server)
	log.Println("card-network running on :60053")
	log.Fatal(server.Serve(lis))
}
