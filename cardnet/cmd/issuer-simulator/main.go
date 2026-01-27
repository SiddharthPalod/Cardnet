package main

import (
	"cardnet/internal/issuer/config"
	"cardnet/internal/issuer/routing"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	issuerpb "cardnet/internal/api/issuer"
	"cardnet/internal/issuer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const issuerConfigPath = "internal/issuer/config/issuers.yaml"

func main() {
	table, err := config.LoadIssuers(issuerConfigPath)
	if err != nil {
		log.Fatal("failed to load issuer config:", err)
	}
	routing.Load(table)
	watchIssuerConfig(issuerConfigPath)

	port := os.Getenv("PORT")
	if port == "" {
		port = "60054"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	issuerpb.RegisterIssuerServiceServer(server, issuer.NewHandler())
	reflection.Register(server)
	log.Println("issuer-simulator running on :60054")
	log.Fatal(server.Serve(lis))
}

func watchIssuerConfig(path string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)

	go func() {
		for range ch {
			log.Println("[issuer-simulator] Reloading Issuer config...")

			table, err := config.LoadIssuers(path)
			if err != nil {
				log.Println("[issuer-simulator] Failed to reload", err)
				continue
			}
			routing.Load(table)
			log.Println("[issuer-simulator] Reloaded Issuer successfully")
		}
	}()
}
