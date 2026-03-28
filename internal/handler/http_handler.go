package handler

import (
	"auth-service/internal/helper"
	"auth-service/internal/middleware"
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHTTPHandler struct {
	Service service.AuthService
}

func NewAuthHTTPHandler(s service.AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{Service: s}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Position string `json:"position"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", nil)
		return
	}

	if err := h.Service.Register(req.Username, req.Password, req.Email, req.Role, req.Position); err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, "Internal Server Error", "Failed to register user: "+err.Error(), nil)
		return
	}

	helper.SendResponse(w, http.StatusCreated, "Created", "User registered successfully", nil)
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendResponse(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", nil)
		return
	}

	claims, token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		helper.SendResponse(w, http.StatusUnauthorized, "Unauthorized", "Invalid credentials", nil)
		return
	}

	middleware.AddUserToLog(r.Context(), claims.UserID, claims.Username, claims.Role)

	helper.SendResponse(w, http.StatusOK, "OK", "Login successful", map[string]string{
		"token": token,
	})
}
