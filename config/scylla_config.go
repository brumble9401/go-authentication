package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ScyllaConfig struct {
	ScyllaHosts         []string `envconfig:"SCYLLA_HOSTS"`
	ScyllaKeyspace      string   `envconfig:"SCYLLA_KEYSPACE"`
	ScyllaMigrationsDir string   `envconfig:"SCYLLA_MIGRATIONS_DIR"`
}

func Load() (ScyllaConfig, error) {
	var c ScyllaConfig

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	err = envconfig.Process("", &c)
	return c, err
}