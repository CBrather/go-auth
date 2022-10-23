package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/CBrather/go-auth/internal/repositories/album"

	"github.com/CBrather/go-auth/internal/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = album.ListByArtist(db, "John Coltrane")
	if err != nil {
		log.Printf("%v", err)
	}

	router := gin.Default()

	api.SetupAlbumRoutes(router, db)

	router.Run("localhost:8080")
}
