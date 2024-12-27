package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	REDIS_ADDR          string 	 `envconfig:"REDIS_ADDR"`
	REDIS_PASSWORD      string   `envconfig:"REDIS_PASSWORD"`
	REDIS_DB 			int		 `envconfig:"REDIS_DB"`
}

func LoadRedisConfig() (RedisConfig, error) {
	var c RedisConfig

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	err = envconfig.Process("", &c)
	return c, err
}