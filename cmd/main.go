package main

import (
	"log"

	"github.com/brumble9401/golang-authentication/api"
	"github.com/brumble9401/golang-authentication/config"
	"github.com/brumble9401/golang-authentication/querybuilder"
	"github.com/brumble9401/golang-authentication/redis"
	"github.com/brumble9401/golang-authentication/repository"
	"github.com/brumble9401/golang-authentication/scylla"
	"github.com/brumble9401/golang-authentication/services/session"
	"github.com/brumble9401/golang-authentication/services/user"
)

func main() {
    cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	manager := scylla.NewManager(cfg)
	if err = manager.CreateKeyspace(cfg.ScyllaKeyspace); err != nil {
		panic(err)
	}
	scyllaSession, err := manager.Connect()
	if err != nil {
		panic(err)
	}
	defer scyllaSession.Close()

	redisCfg, err := config.LoadRedisConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("Connecting to Redis at %s", redisCfg.REDIS_ADDR)
	redisClient := redis.NewRedisClient(redisCfg.REDIS_ADDR, redisCfg.REDIS_PASSWORD, redisCfg.REDIS_DB)
	redisService := session.NewService(redisClient)

    queryBuilder := querybuilder.NewScyllaQueryBuilder(scyllaSession.Session)
	userRepo := repository.NewUserRepository(scyllaSession.Session, queryBuilder)
	roleRepo := repository.NewRoleRepository(scyllaSession.Session, queryBuilder)
	authProviderRepo := repository.NewAuthProviderRepository(scyllaSession.Session, queryBuilder, *redisService)
    userService := user.NewService(userRepo, roleRepo, authProviderRepo, *redisService)

    apiServer := api.NewAPIServer(":8080", userService)
    log.Fatal(apiServer.Run())
}