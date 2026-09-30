package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"service/internal/database"
	"service/internal/handler"
	"service/internal/middleware"
	"service/internal/repository"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://domservice:domservice_local_dev@127.0.0.1:5434/domservice?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer pool.Close()

	if err := database.CreateTables(ctx, pool); err != nil {
		log.Fatal("could not create tables: ", err)
	}

	log.Println("PostgreSQL connected successfully")

	repo := repository.NewRequestRepository(pool)
	requestHandler := handler.NewRequestHandler(repo)

	userRepo := repository.NewUserRepository(pool)
	authHandler := handler.NewAuthHandler(userRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":   "ok",
			"service":  "domservice",
			"database": "connected",
		})
	})

	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	mux.HandleFunc("POST /api/requests", requestHandler.Create)

	mux.Handle(
		"GET /api/requests",
		middleware.Auth(
			middleware.AdminOnly(
				http.HandlerFunc(requestHandler.List),
			),
		),
	)

	mux.Handle(
		"PATCH /api/requests/{id}/status",
		middleware.Auth(
			middleware.AdminOnly(
				http.HandlerFunc(requestHandler.UpdateStatus),
			),
		),
	)

	mux.Handle(
		"DELETE /api/requests/{id}",
		middleware.Auth(
			middleware.AdminOnly(
				http.HandlerFunc(requestHandler.Delete),
			),
		),
	)

	// Small frontend.
	fileServer := http.FileServer(http.Dir("./web"))
	mux.Handle("GET /", fileServer)

	server := &http.Server{
		Addr:              "127.0.0.1:" + port,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("DomService started: http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}

	log.Println("server stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
