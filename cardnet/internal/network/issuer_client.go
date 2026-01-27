package network

import (
	issuerpb "cardnet/internal/api/issuer"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewIssuerClient() issuerpb.IssuerServiceClient {
	// Connect to the monolithic issuer simulator
	addr := os.Getenv("ISSUER_SIMULATOR_ADDR")
	if addr == "" {
		addr = "localhost:60054"
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return issuerpb.NewIssuerServiceClient(conn)
}
