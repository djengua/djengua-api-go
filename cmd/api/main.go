package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mongoadapter "github.com/djengua/djengua-api-go/internal/adapters/mongo"
	"github.com/djengua/djengua-api-go/internal/application"
	"github.com/djengua/djengua-api-go/internal/config"
	"github.com/djengua/djengua-api-go/internal/httpapi/middleware"
	"github.com/djengua/djengua-api-go/internal/httpapi/router"
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
	productsService := application.NewProductsService(repo)
	categoriesService := application.NewCategoriesService(repo)
	collectionsService := application.NewCollectionsService(repo)

	log.Printf("CORS ALLOWED: %v", cfg.AllowedOrigins)
	h := middleware.CORS(router.New(productsService, categoriesService, collectionsService), cfg.AllowedOrigins)

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
