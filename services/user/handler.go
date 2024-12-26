package user

import (
	"encoding/json"
	"net/http"

	"github.com/brumble9401/golang-authentication/interfaces"
	"github.com/brumble9401/golang-authentication/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/mux"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/register", h.Register).Methods("POST")
    router.HandleFunc("/login", h.Login).Methods("POST")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    log.Info("Registering user")
	var payload models.RegisterUserPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    if err := h.service.Register(r.Context(), &payload); err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := interfaces.NewResponse()
    response.SetStatus("success")
    response.SetMessage("User registered successfully")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    log.Info("Logging in user")
    var payload models.LoginPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    token, err := h.service.Login(r.Context(), &payload)
    if err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := interfaces.NewResponse()
    response.SetStatus("success")
    response.SetMessage("Login successful")
    response.SetData(map[string]string{"token": token})
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}