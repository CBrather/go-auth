package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	adapter "github.com/gwatts/gin-adapter"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api/middleware"
	"github.com/CBrather/go-auth/internal/repositories/album"
)

var db *sql.DB

func SetupAlbumRoutes(rootRouter *chi.Mux, newDb *sql.DB) {
	db = newDb

	albumRouter := chi.NewRouter()
	albumRouter.Use(middleware.EnsureValidToken())

	albumRouter.With(middleware.RequireScope("read:albums")).Get("/", listAlbums)

	rootRouter.Mount("/albums", albumRouter)
}

func listAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := album.List(db)
	if err != nil {
		zap.L().Error("Failed to retrieve the list of albums.", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zap.L().Info("Successfully retrieved a list of albums to return")

	body, err := json.Marshal(albums)
	if err != nil {
		zap.L().Error("Failed to serialize list of albums", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

/*
** Gin Starts Here
 */

func SetupAlbumRoutesGin(router *gin.Engine, setupDb *sql.DB) {
	db = setupDb

	validateToken := adapter.Wrap(middleware.EnsureValidToken())
	router.GET("/album/:id", validateToken, middleware.RequireScopeGin("read:albums"), getAlbumGin)
	router.POST("/album", validateToken, middleware.RequireScopeGin("create:albums"), postAlbumGin)
}

func getAlbumGin(ginCtx *gin.Context) {
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

func postAlbumGin(ginCtx *gin.Context) {
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
