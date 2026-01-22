// internal/core/usecase/auth/service.go
package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/djengua/djengua-api-go/internal/core/domain"
	"github.com/djengua/djengua-api-go/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users     ports.AuthUserRepository
	jwtSecret []byte
	issuer    string
	ttl       time.Duration
}

func NewService(users ports.AuthUserRepository, jwtSecret string, issuer string, ttl time.Duration) *Service {
	return &Service{
		users:     users,
		jwtSecret: []byte(jwtSecret),
		issuer:    issuer,
		ttl:       ttl,
	}
}

type Claims struct {
	UserID string          `json:"uid"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) Register(ctx context.Context, email, password, name string) (domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	if email == "" || password == "" || name == "" {
		return domain.User{}, "", errors.New("email, password, name are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, "", err
	}

	u, err := s.users.CreateUser(ctx, domain.User{
		Email:        email,
		Name:         name,
		Role:         domain.RoleUser,
		PasswordHash: string(hash),
	})
	if err != nil {
		return domain.User{}, "", err
	}

	token, err := s.NewToken(u)
	if err != nil {
		return domain.User{}, "", err
	}

	// Nunca exponer hash
	u.PasswordHash = ""
	return u, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return domain.User{}, "", errors.New("email and password are required")
	}

	u, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		// No revelar si existe o no
		return domain.User{}, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, "", errors.New("invalid credentials")
	}

	token, err := s.NewToken(u)
	if err != nil {
		return domain.User{}, "", err
	}

	u.PasswordHash = ""
	return u, token, nil
}

func (s *Service) NewToken(u domain.User) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   u.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}

func (s *Service) ParseToken(tokenString string) (ports.AuthClaims, error) {
	tok, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return s.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return ports.AuthClaims{}, err
	}

	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return ports.AuthClaims{}, jwt.ErrTokenInvalidClaims
	}

	// Validaciones extra opcionales (issuer)
	if s.issuer != "" && claims.Issuer != s.issuer {
		return ports.AuthClaims{}, jwt.ErrTokenInvalidClaims
	}

	return ports.AuthClaims{UserID: claims.UserID, Role: claims.Role}, nil
}

func (s *Service) GetUser(ctx context.Context, id string) (domain.User, error) {
	u, err := s.users.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u.PasswordHash = ""
	return u, nil
}
