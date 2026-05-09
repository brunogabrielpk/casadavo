package handler

import (
	"net/http"

	mw "casadavo/internal/middleware"
	"casadavo/internal/repository"
	"casadavo/internal/service"
)

type AuthHandler struct {
	svc   *service.AuthService
	users *repository.UserRepo
}

func NewAuthHandler(svc *service.AuthService, users *repository.UserRepo) *AuthHandler {
	return &AuthHandler{svc, users}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := readJSON(r, &body); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.Name == "" || body.Email == "" || body.Password == "" {
		errResp(w, http.StatusBadRequest, "name, email and password are required")
		return
	}
	u, err := h.svc.Register(body.Name, body.Email, body.Phone, body.Password, body.Role)
	if err != nil {
		errResp(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &body); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	token, u, err := h.svc.Login(body.Email, body.Password)
	if err != nil {
		errResp(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	c := mw.GetClaims(r)
	u, err := h.users.FindByID(c.UserID)
	if err != nil || u == nil {
		errResp(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, u)
}
