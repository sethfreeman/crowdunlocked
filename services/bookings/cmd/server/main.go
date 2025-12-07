package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/crowdunlocked/services/bookings/internal/handler"
	"github.com/crowdunlocked/services/bookings/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Configure DynamoDB client with optional local endpoint
	var dynamoClient *dynamodb.Client
	if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
		dynamoClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &endpoint
		})
		log.Printf("Using DynamoDB endpoint: %s", endpoint)
	} else {
		dynamoClient = dynamodb.NewFromConfig(cfg)
	}

	tableName := getEnv("DYNAMODB_TABLE", "bookings")
	repo := repository.NewBookingRepository(dynamoClient, tableName)
	bookingHandler := handler.NewBookingHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/bookings", func(r chi.Router) {
		r.Post("/", bookingHandler.Create)
		r.Get("/{id}", bookingHandler.GetByID)
		r.Post("/{id}/confirm", bookingHandler.Confirm)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing health check response: %v", err)
		}
	})

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting bookings service on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
