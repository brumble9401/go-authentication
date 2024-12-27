package api

import (
	"net/http"

	"github.com/brumble9401/golang-authentication/services/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
    addr string
    userService user.Service
}

func NewAPIServer(addr string, userService user.Service) *APIServer {
    return &APIServer{
        addr: addr,
        userService: userService,
    }
}

func (s *APIServer) Run() error {
    router := mux.NewRouter()
    subRouter := router.PathPrefix("/api/v1").Subrouter()
    // subRouter.Use(middleware.AuthMiddleware)

    userHandler := user.NewHandler(s.userService)
    userHandler.RegisterRoutes(subRouter)

    logrus.Info("Listening on ", s.addr)

    return http.ListenAndServe(s.addr, router)
}