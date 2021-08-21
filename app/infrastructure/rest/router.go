package rest

import (
	"geo/infrastructure/env"
	"geo/interfaces/rest/handler"
	gintraceser "github.com/AeroAgency/go-gin-tracer"
	"github.com/AeroAgency/golang-helpers-lib/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	// Обработчик запросов geo
	geoHandler *handler.GeoHandler
}

// Конструктор
func NewRouter(geoHandler *handler.GeoHandler) *Router {
	return &Router{
		geoHandler: geoHandler,
	}
}

func (p Router) Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gintraceser.OpenTracingMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
	}))

	eapi := r.Group(env.Root)
	v := eapi.Group(env.Version)
	service := v.Group("/" + env.ServiceName)
	service.Use(logger.Logger())
	//manage methods
	manage := service.Group("/manage")
	manage.GET("/health", handler.Health)
	//api methods
	service.GET("/detectAddressInfo", p.geoHandler.DetectAddressInfo)
	service.POST("/suggest/address", p.geoHandler.SuggestAddress)
	service.GET("/findById/address", p.geoHandler.FindAddressById)
	return r
}
