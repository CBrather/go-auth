package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api"
	"github.com/CBrather/go-auth/pkg/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// Only logging on Warn level, as this might well be intentional, especially in production environments
		zap.L().Warn("Unable to load a .env file, will execute with environment as-is", zap.Error(err))
	}

	if err := log.Initialize(os.Getenv("LOGLEVEL")); err != nil {
		zap.L().Fatal("Failed to setup logger")
	} else {
		zap.L().Info("Logger was successfully setup")
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE")))
	if err != nil {
		zap.L().Fatal("Unable to open a Postgres connection", zap.Error(err))
	}

	setupRoutes(db)
}

func setupRoutes(db *sql.DB) {
	logger := httplog.NewLogger("go-auth", httplog.Options{JSON: true, Concise: true})
	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.Recoverer)

	router.Handle("/metrics", promhttp.Handler())
	api.SetupProbeRoutes(router)

	api.SetupAlbumRoutes(router, db)

	zap.L().Info("Server listening on :8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}
