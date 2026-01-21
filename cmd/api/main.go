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
	"github.com/djengua/djengua-api-go/internal/core/usecase/categories"
	"github.com/djengua/djengua-api-go/internal/core/usecase/collections"
	"github.com/djengua/djengua-api-go/internal/core/usecase/products"
	"github.com/djengua/djengua-api-go/internal/platform/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.FromEnv()

	ctx := context.Background()
	client, err := mongoadapter.NewClient(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connection failed: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database(cfg.MongoDB)
	repo := mongoadapter.NewRepository(database)

	productsUC := products.NewService(repo)
	categoriesUC := categories.NewService(repo)
	collectionsUC := collections.NewService(repo)

	log.Printf("CORS ALLOWED: %v", cfg.AllowedOrigins)
	h := httpmw.CORS(router.New(productsUC, categoriesUC, collectionsUC), cfg.AllowedOrigins)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
