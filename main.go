package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/albums", listAlbums)
	router.GET("/albums/:id", getAlbum)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func listAlbums(ginCtx *gin.Context) {
	ginCtx.IndentedJSON(http.StatusOK, albums)
}

func getAlbum(ginCtx *gin.Context) {
	id := ginCtx.Param("id")

	for _, album := range albums {
		if album.ID == id {
			ginCtx.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	ginCtx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbums(ginCtx *gin.Context) {
	var newAlbum album

	if err := ginCtx.BindJSON(&newAlbum); err != nil {
		fmt.Printf("Deserializing request failed: %v", err)
		return
	}

	albums = append(albums, newAlbum)

	ginCtx.IndentedJSON(http.StatusCreated, newAlbum)
}
