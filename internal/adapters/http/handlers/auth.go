// internal/adapters/http/handlers/auth.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/core/ports"
	"github.com/djengua/djengua-api-go/internal/core/usecase/auth"
)

type AuthHandler struct {
	auth  *auth.Service
	users ports.AuthUserRepository
}

func NewAuthHandler(authSvc *auth.Service, users ports.AuthUserRepository) *AuthHandler {
	return &AuthHandler{auth: authSvc, users: users}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in registerReq
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(in.Email) == "" || strings.TrimSpace(in.Password) == "" || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "email, password, name are required")
		return
	}

	u, tok, err := h.auth.Register(r.Context(), in.Email, in.Password, in.Name)
	if err != nil {
		if err == ports.ErrConflict {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, authResp{Token: tok, User: u})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in loginReq
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	u, tok, err := h.auth.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	writeJSON(w, http.StatusOK, authResp{Token: tok, User: u})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing auth")
		return
	}
	u, err := h.users.GetUserByID(r.Context(), uid)
	if mapDBErr(w, err) {
		return
	}
	u.PasswordHash = ""
	writeJSON(w, http.StatusOK, u)
}
