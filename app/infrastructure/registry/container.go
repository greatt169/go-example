package registry

import (
	"geo/application/service"
	"geo/infrastructure/env"
	"geo/infrastructure/interfaces/geo"
	infrDadata "geo/infrastructure/persistence/dadata"
	"geo/infrastructure/rest"
	"geo/interfaces/rest/handler"
	redisCache "github.com/AeroAgency/redis-cache"
	"github.com/gomodule/redigo/redis"
	zerologger "github.com/rs/zerolog/log"
	"github.com/sarulabs/di"
	log "github.com/sirupsen/logrus"
	"gopkg.in/webdeskltd/dadata.v2"
	"os"
)

// Контейнер внедрения зависимостей
type Container struct {
	ctn di.Container
}

// Конструктор для контейнера
func NewContainer() (*Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}
	if err := builder.Add([]di.Def{
		{ // router
			Name:  "router",
			Build: buildRouterContainer,
		},
		{ // Сервис работы с ошибками
			Name:  "error_service",
			Build: buildErrorServiceContainer,
		},
		{ // Клиент для работы с dadata api
			Name:  "dadata_client",
			Build: buildDadataClientContainer,
		},
		{ // Адаптер для dadata api
			Name:  "geo_adapter",
			Build: buildGeoAdapterContainer,
		},
		{ // Сервис для работы с геопозициями
			Name:  "geo_service",
			Build: buildGeoServiceContainer,
		},
		{ // Обработчик запросов к geo
			Name:  "geo_handler",
			Build: buildGeoHandlerContainer,
		},
		{ // Логгер
			Name: "logger",
			Build: func(ctn di.Container) (interface{}, error) {
				logger := log.Logger{}
				logger.SetFormatter(&log.JSONFormatter{})
				logger.SetOutput(os.Stdout)
				logger.Info("run rest server")
				logger.SetLevel(log.DebugLevel)
				return &logger, nil
			},
		},
		{ // Redis
			Name:  "redis",
			Build: buildRedisContainer,
		},
		{ // Сервис работы с кешированием
			Name:  "cache_service",
			Build: buildCacheServiceContainer,
		},
	}...); err != nil {
		return nil, err
	}

	return &Container{
		ctn: builder.Build(),
	}, nil
}

// Собирает контейнер с Redis
func buildRedisContainer(ctn di.Container) (interface{}, error) {
	c, err := redis.Dial("tcp", env.RedisHost+":"+env.RedisPort)
	if err != nil {
		log.Fatal(err)
	}
	return c, err
}

// собирает контейнер для работы с ошибками
func buildErrorServiceContainer(ctn di.Container) (interface{}, error) {
	// Получение лог-сервиса из контейнера
	logger := ctn.Get("logger").(log.FieldLogger)
	errorService := service.NewError(logger)
	return errorService, nil
}

// собирает контейнер для работы с клиентом dadata
func buildDadataClientContainer(ctn di.Container) (interface{}, error) {
	daData := dadata.NewDaData(env.DadataApiKey, env.DadataSecretKey)
	return daData, nil
}

// собирает контейнер адаптера для работы с api dadata
func buildGeoAdapterContainer(ctn di.Container) (interface{}, error) {
	daDataClient := ctn.Get("dadata_client").(*dadata.DaData)
	geoGeoAdapter := infrDadata.NewDadataGeoAdapter(daDataClient)
	return geoGeoAdapter, nil
}

// собирает контейнер для работы с ошибками
func buildGeoServiceContainer(ctn di.Container) (interface{}, error) {
	geoGeoAdapter := ctn.Get("geo_adapter").(geo.GeoAdapterInterface)
	geoService := service.NewGeoService(geoGeoAdapter)
	return geoService, nil
}

//Собирает контейнер для работы с хендлером geo
func buildGeoHandlerContainer(ctn di.Container) (interface{}, error) {
	// Получение сервиса для работы с ошибками
	errorService := ctn.Get("error_service").(*service.Error)
	// Получение сервиса для работы с геолокациями
	geoService := ctn.Get("geo_service").(*service.GeoService)
	// Получение сервиса для работы с кешем
	cacheService := ctn.Get("cache_service").(*redisCache.CacheService)
	server := handler.NewGeoHandler(errorService, geoService, cacheService)
	return server, nil
}

//Собирает контейнер для роутера
func buildRouterContainer(ctn di.Container) (interface{}, error) {
	geoHandler := ctn.Get("geo_handler").(*handler.GeoHandler)
	server := rest.NewRouter(geoHandler)
	return server, nil
}

//Собирает контейнер для работы с сервисом кеширования
func buildCacheServiceContainer(ctn di.Container) (interface{}, error) {
	redisContainer := ctn.Get("redis").(redis.Conn)
	cacheService := redisCache.NewCacheService(zerologger.Logger, redisContainer)
	return cacheService, nil
}

// Получение зависимости из контейнера
func (c *Container) Resolve(name string) interface{} {
	return c.ctn.Get(name)
}

// Очистка контейнера
func (c *Container) Clean() error {
	return c.ctn.Clean()
}
