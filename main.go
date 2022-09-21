package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	albumAPI "go-auth/internal/api/album"
	albumModel "go-auth/internal/repository/album"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = albumModel.ListByArtist(db, "John Coltrane")
	if err != nil {
		log.Printf("%v", err)
	}

	router := gin.Default()

	albumAPI.SetupRoutes(router)

	router.Run("localhost:8080")
}
