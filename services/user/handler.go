package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brumble9401/golang-authentication/config"
	"github.com/brumble9401/golang-authentication/interfaces"
	"github.com/brumble9401/golang-authentication/middleware"
	"github.com/brumble9401/golang-authentication/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type Handler struct {
    service Service
    googleOauthConfig oauth2.Config
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service, googleOauthConfig: config.LoadGoogleConfig()}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/register", h.Register).Methods("POST")
    router.HandleFunc("/login", h.Login).Methods("POST")
    // router.HandleFunc("/update-username-password", h.UpdateUsernameAndPassword).Methods("PUT").Subrouter().Use(middleware.AuthMiddleware)
    router.HandleFunc("/login/google", h.GoogleLogin).Methods("GET")
    router.HandleFunc("/callback", h.GoogleCallback).Methods("GET")

    router.Handle("/update-username-password", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateUsernameAndPassword))).Methods("PUT")

    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}), // Adjust the allowed origins as needed
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )

    http.Handle("/", corsHandler(router))
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

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
    log.Info("GoogleLogin handler invoked")
    b := make([]byte, 16)
    _, _ = rand.Read(b)
    state := base64.URLEncoding.EncodeToString(b)
    url := h.googleOauthConfig.AuthCodeURL(state)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
    log.Info("GoogleCallback handler invoked")

    code := r.FormValue("code")
    token, err := h.googleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
        return
    }

    client := h.googleOauthConfig.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
        return
    }

    accessToken, err := h.service.LoginByGoogle(r.Context(), userInfo)
    if err != nil {
        http.Error(w, "Failed to login with Google: "+err.Error(), http.StatusInternalServerError)
        return
    }

    redirectURL := fmt.Sprintf("http://localhost:5174?access_token=%s", accessToken)
    http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *Handler) UpdateUsernameAndPassword(w http.ResponseWriter, r *http.Request) {
    log.Info("Updating username and password")
    log.Debug("User id in context: ", r.Context().Value("userID"))
    var payload models.UsernamePasswordPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    if err := h.service.UpdateUsernameAndPassword(r.Context(), &payload); err != nil {
        response := interfaces.NewResponse()
        response.SetStatus("error")
        response.SetMessage(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := interfaces.NewResponse()
    response.SetStatus("success")
    response.SetMessage("Username and password updated successfully")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}