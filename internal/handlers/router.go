package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/serjnester/weather-subscription-service/internal/configs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type RouterParams struct {
	Config configs.Config
}

// NewRouter godoc
//
//	@title			Weather Forecast API
//	@version		1.0
//	@description	Weather API application that allows users to subscribe to weather updates for their city.
func NewRouter(params RouterParams) *gin.Engine {
	if !params.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	router.NoMethod(missingRoutePath)
	router.NoRoute(missingRoutePath)

	router.GET("/health/readiness", ReadinessProbe)
	router.GET("/health/liveness", LivenessProbe)

	if !params.Config.Env.IsProd() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}

func missingRoutePath(c *gin.Context) {
	c.AbortWithError(http.StatusNotFound, fmt.Errorf("path: [%s: %s] is missing", c.Request.Method, c.Request.URL.Path))
}

type RegisterHandlersParams struct {
	MainHandler Handler
}

func RegisterHandlers(router *gin.Engine, params RegisterHandlersParams) {
	handler := params.MainHandler

	api := router.Group("/api")
	api.GET("/weather", handler.GetWeather)
	api.POST("/subscribe", handler.Subscribe)
	api.GET("/confirm/:token", handler.ConfirmSubscription)
	api.GET("/unsubscribe/:token", handler.Unsubscribe)
}
