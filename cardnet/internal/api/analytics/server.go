package analytics

import (
	"cardnet/internal/analytics/repository/postgres"
)

type Server struct {
	UnimplementedAnalyticsServiceServer
	repo *postgres.Repository
}

func NewServer(repo *postgres.Repository) *Server {
	return &Server{
		repo: repo,
	}
}
