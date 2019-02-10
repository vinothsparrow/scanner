package server

import (
	"github.com/gin-gonic/gin"
	"github.com/vinothsparrow/scanner/config"
	"github.com/vinothsparrow/scanner/controllers"
	"github.com/vinothsparrow/scanner/middlewares"
)

func NewRouter() *gin.Engine {
	config := config.GetConfig()
	gin.SetMode(config.GetString("http.env"))
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.AuthMiddleware())
	router.Use(middlewares.ErrorMiddleware())

	notfound := new(controllers.NotFoundController)
	router.NoRoute(notfound.Error404)
	router.NoMethod(notfound.Error405)
	v1 := router.Group("v1")
	{
		scanGroup := v1.Group("scan")
		{
			scan := new(controllers.ScanController)
			scanGroup.GET("/git", scan.GitScan)
		}
	}
	return router

}
