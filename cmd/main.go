package main

import (
	"log"

	"github.com/brumble9401/golang-authentication/api"
	"github.com/brumble9401/golang-authentication/config"
	"github.com/brumble9401/golang-authentication/querybuilder"
	"github.com/brumble9401/golang-authentication/repository"
	"github.com/brumble9401/golang-authentication/scylla"
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

	session, err := manager.Connect()
	if err != nil {
		panic(err)
	}

	defer session.Close()

    queryBuilder := querybuilder.NewScyllaQueryBuilder(session.Session)
	userRepo := repository.NewUserRepository(session.Session, queryBuilder)
	roleRepo := repository.NewRoleRepository(session.Session, queryBuilder)
    userService := user.NewService(userRepo, roleRepo)

    apiServer := api.NewAPIServer(":8080", userService)
    log.Fatal(apiServer.Run())
}