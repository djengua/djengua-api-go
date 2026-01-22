// internal/adapters/http/middleware/jwt.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/djengua/djengua-api-go/internal/core/ports"
)

type ctxKey string

const (
	CtxUserID ctxKey = "uid"
	CtxRole   ctxKey = "role"
)

func RequireJWT(authSvc ports.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
				http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
				return
			}
			raw := strings.TrimSpace(h[len("Bearer "):])

			claims, err := authSvc.ParseToken(raw)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	v := r.Context().Value(CtxUserID)
	s, ok := v.(string)
	return s, ok && s != ""
}
