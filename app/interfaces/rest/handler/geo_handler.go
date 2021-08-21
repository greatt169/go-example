package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"geo/application/dto"
	"geo/application/service"
	"geo/infrastructure/env"
	appErrors "geo/infrastructure/errors"
	redisCache "github.com/AeroAgency/redis-cache"
	"github.com/gin-gonic/gin"
	pkgErrors "github.com/pkg/errors"
	"net/http"
	"strconv"
)

// GeoHandler Обработчик для запросов сервиса местоположений
type GeoHandler struct {
	errorService     *service.Error
	geoService       *service.GeoService
	validatorService *service.Validator
	cacheService     *redisCache.CacheService
	redisExpTime     int
}

// NewGeoHandler Конструктор
func NewGeoHandler(errorService *service.Error, geoService *service.GeoService, cacheService *redisCache.CacheService) *GeoHandler {
	redisExpTime, _ := strconv.Atoi(env.RedisExpTime)
	return &GeoHandler{
		errorService:     errorService,
		geoService:       geoService,
		validatorService: &service.Validator{},
		cacheService:     cacheService,
		redisExpTime:     redisExpTime,
	}
}

// DetectAddressInfo
// Определение информации о местоположении пользователя
func (gh *GeoHandler) DetectAddressInfo(c *gin.Context) {
	latString := c.Query("lat")
	lonString := c.Query("lon")
	// проверка
	if latString == "" || lonString == "" {
		gh.errorService.HandleError(appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("empty value of params [lat, lon] is invalid"))}, c)
		return
	}
	var lat float64
	var lon float64
	lat, _ = strconv.ParseFloat(latString, 64)
	lon, _ = strconv.ParseFloat(lonString, 64)
	if lat == 0 {
		gh.errorService.HandleError(appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("incorrect val of param lat"))}, c)
		return
	}
	if lon == 0 {
		gh.errorService.HandleError(appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("incorrect val ofparam lon"))}, c)
	}
	// Кэширование
	tag := fmt.Sprintf("cache:detect_address_info_%v_%v", lat, lon)
	locationResponse := &dto.Location{}
	result := gh.cacheService.GetByTag(tag, locationResponse)
	if result == false {
		var err error
		// Вызов сервиса
		locationResponse, err = gh.geoService.DetectAddressInfoRecursive(lat, lon)
		if err != nil {
			gh.errorService.HandleError(err, c)
			return
		}
		gh.cacheService.SetByTag(tag, locationResponse, gh.redisExpTime)
	}
	c.JSON(http.StatusOK, locationResponse)
}

// SuggestAddress
// Возвращает подсказки по местоположениям
func (gh *GeoHandler) SuggestAddress(c *gin.Context) {
	var suggestParams dto.SuggestsRequestDto
	err := c.BindJSON(&suggestParams)
	if err != nil {
		gh.errorService.HandleError(err, c)
		return
	}
	// Валидируем основные поля
	err = gh.validatorService.ValidateSuggestAddressDto(suggestParams)
	if err != nil {
		gh.errorService.HandleError(err, c)
		return
	}
	// Преобразуем параметры в json
	b, err := json.Marshal(suggestParams)
	if err != nil {
		gh.errorService.HandleError(err, c)
		return
	}
	// Преобразуем параметры в base64
	tag := base64.StdEncoding.EncodeToString(b)
	suggests := &[]dto.LocationAddress{}
	result := gh.cacheService.GetByTag(tag, suggests)
	if result == false {
		// Вызов сервиса
		suggests, err = gh.geoService.SuggestAddress(suggestParams)
		if err != nil {
			gh.errorService.HandleError(err, c)
			return
		}
		gh.cacheService.SetByTag(tag, suggests, gh.redisExpTime)
	}
	c.JSON(http.StatusOK, suggests)
}

// FindAddressById
// Находит адрес по идентификатору: ФИСА-код/КЛАДР-код
func (gh *GeoHandler) FindAddressById(c *gin.Context) {
	var findAddressByIdParams dto.FindAddressByIdRequestDto
	findAddressByIdParams.Query = c.Query("query")
	// Валидация
	err := gh.validatorService.ValidateFindAddressByIdDto(findAddressByIdParams)
	if err != nil {
		gh.errorService.HandleError(err, c)
		return
	}
	// Кэширование
	tag := fmt.Sprintf("cache:find_address_by_id_%s", c.Query("query"))
	address := &dto.LocationAddress{}
	result := gh.cacheService.GetByTag(tag, address)
	if result == false {
		var err error
		// Вызов сервиса
		address, err = gh.geoService.FindAddressById(findAddressByIdParams)
		if err != nil {
			gh.errorService.HandleError(err, c)
			return
		}
		gh.cacheService.SetByTag(tag, address, gh.redisExpTime)
	}
	c.JSON(http.StatusOK, address)
}
