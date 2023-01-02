package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api/middleware"
	"github.com/CBrather/go-auth/internal/repositories/album"
)

var db *sql.DB

func SetupAlbumRoutes(router *gin.Engine, setupDb *sql.DB) {
	db = setupDb

	validateToken := adapter.Wrap(middleware.EnsureValidToken())
	router.GET("/album/:id", validateToken, middleware.RequireScope("read:albums"), getAlbum)
	router.GET("/album", validateToken, middleware.RequireScope("read:albums"), listAlbums)
	router.POST("/album", validateToken, middleware.RequireScope("create:albums"), postAlbum)
}

func listAlbums(ginCtx *gin.Context) {
	albums, err := album.List(db)
	if err != nil {
		zap.L().Error("Failed to retrieve the list of albums.", zap.Error(err))

		ginCtx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.L().Info("Successfully retrieved a list of albums to return")

	ginCtx.IndentedJSON(http.StatusOK, albums)
}

func getAlbum(ginCtx *gin.Context) {
	idString := ginCtx.Params.ByName("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		ginCtx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	id64 := int64(id)
	album, err := album.GetByID(db, id64)

	if err != nil {
		ginCtx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

		return
	}

	ginCtx.IndentedJSON(http.StatusOK, album)
}

func postAlbum(ginCtx *gin.Context) {
	var newAlbum album.Album
	if err := ginCtx.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := album.Add(db, newAlbum)
	if err != nil {
		ginCtx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed saving the requested album"})
	}

	ginCtx.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}
