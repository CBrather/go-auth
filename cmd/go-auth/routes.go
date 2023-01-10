package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api"
	"github.com/CBrather/go-auth/pkg/telemetry"
)

func SetupHttpRoutes(db *sql.DB) {
	logger := httplog.NewLogger("go-auth", httplog.Options{JSON: true, Concise: true})
	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.Recoverer)

	router.Handle("/metrics", telemetry.NewMetricsHandler())
	api.SetupProbeRoutes(router)

	api.SetupAlbumRoutes(router, db)

	zap.L().Info("Server listening on :8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}
