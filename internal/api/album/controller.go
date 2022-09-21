package album

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/albums", listAlbums)
	router.GET("/albums/:id", getAlbum)
	router.POST("/albums", postAlbums)
}

func listAlbums(ginCtx *gin.Context) {
	ginCtx.IndentedJSON(http.StatusOK, make([]byte, 1)) //TODO: Hook up with model
}

func getAlbum(ginCtx *gin.Context) { //TODO: Hook up with model
	ginCtx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbums(ginCtx *gin.Context) {
	/*
		var newAlbum Album
		if err := ginCtx.BindJSON(&newAlbum); err != nil {
			fmt.Printf("Deserializing request failed: %v", err)
			return
		}
	*/
	ginCtx.IndentedJSON(http.StatusCreated, make([]byte, 1)) //TODO: Hook up with model
}
