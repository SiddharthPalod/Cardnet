package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

type Config struct {
	Hosts       []string
	Keyspace    string
	Consistency gocql.Consistency
	Timeout     time.Duration
}

func NewSession(cfg Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)

	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = cfg.Consistency
	cluster.Timeout = cfg.Timeout

	// ---- Production defaults ----
	cluster.NumConns = 4
	cluster.ReconnectInterval = 5 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}

	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true

	// Observability (minimal but useful)
	cluster.Events.DisableTopologyEvents = true
	cluster.Events.DisableSchemaEvents = true

	log.Printf(
		"connecting to cassandra hosts=%v keyspace=%s",
		cfg.Hosts, cfg.Keyspace,
	)

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
