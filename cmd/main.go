package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/CBrather/go-auth/internal/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading the .env file: %v, will execute with environment as-is", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE")))
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	api.SetupAlbumRoutes(router, db)
	api.SetupProbeRoutes(router)

	router.Run("0.0.0.0:8080")
}
