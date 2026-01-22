// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpmw "github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/adapters/http/router"
	mongoadapter "github.com/djengua/djengua-api-go/internal/adapters/mongo"
	"github.com/djengua/djengua-api-go/internal/core/usecase/auth"
	"github.com/djengua/djengua-api-go/internal/core/usecase/categories"
	"github.com/djengua/djengua-api-go/internal/core/usecase/collections"
	"github.com/djengua/djengua-api-go/internal/core/usecase/orders"
	"github.com/djengua/djengua-api-go/internal/core/usecase/products"
	"github.com/djengua/djengua-api-go/internal/core/usecase/sales"
	"github.com/djengua/djengua-api-go/internal/platform/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.FromEnv()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Mongo connect con timeout
	mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelMongo()

	client, err := mongoadapter.NewClient(mongoCtx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connection failed: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database(cfg.MongoDB)
	repo := mongoadapter.NewRepository(database)

	// Indexes (ej. users.email unique)
	idxCtx, cancelIdx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelIdx()
	if err := repo.EnsureIndexes(idxCtx); err != nil {
		log.Printf("warning: ensure indexes: %v", err)
	}

	// Usecases
	productsUC := products.NewService(repo)
	categoriesUC := categories.NewService(repo)
	collectionsUC := collections.NewService(repo)
	ordersUC := orders.NewService(repo)
	salesUC := sales.NewService(repo, repo)

	// ✅ AuthUserRepository (mínimo para auth/me)
	// repo debe implementar: CreateUser, GetUserByID, GetUserByEmail
	usersRepo := repo

	// Auth usecase
	authUC := auth.NewService(
		usersRepo,
		cfg.JWTSecret,
		cfg.JWTIssuer,
		time.Duration(cfg.JWTTTLMinutes)*time.Minute,
	)

	log.Printf("CORS ALLOWED: %v", cfg.AllowedOrigins)

	// Router con auth + middleware JWT en writes
	h := httpmw.CORS(
		router.New(productsUC, categoriesUC, collectionsUC, authUC, usersRepo, ordersUC, salesUC),
		cfg.AllowedOrigins,
	)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	log.Printf("shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
