package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupProbeRoutes(router *gin.Engine) {
	router.GET("/healthz", getHealth)
}

func getHealth(ginCtx *gin.Context) {
	ginCtx.AbortWithStatus(http.StatusOK)
}
