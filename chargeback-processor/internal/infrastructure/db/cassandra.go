package db

import (
	"github.com/gocql/gocql"
	"log"
	"processor/internal/config"
	"time"
)

func NewCassandraSession(cfg config.DatabaseConfig) *gocql.Session {
	cluster := gocql.NewCluster(cfg.CassandraHosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	log.Println("Connected to Cassandra!")
	return session
}
