package service

import (
	"fmt"
	"geo/application/dto"``
	appErrors "geo/infrastructure/errors"
	"geo/infrastructure/interfaces/geo"
	pkgErrors "github.com/pkg/errors"
	"strconv"
)

// Широта центра МСК
const defaultLat = 55.7522200

// Долгота центра МСК
const defaultLon = 37.6155600

var fiasLevelMap = map[string]string{
	"0":  "country",        // страна
	"1":  "region",         // регион
	"3":  "area",           // район
	"4":  "city",           // город
	"5":  "city_district",  // район города
	"6":  "settlement",     // населенный пункт
	"7":  "street",         // улица
	"8":  "house",          // дом
	"9":  "flat",           // квартира
	"65": "plan_structure", // планировочная структура
	"-1": "unknown",        // иностранный или пустой
}

// GeoService Сервис для работы с геолопозициями
type GeoService struct {
	geoAdapter geo.GeoAdapterInterface
}

// NewGeoService Конструктор
func NewGeoService(
	geoAdapter geo.GeoAdapterInterface,
) *GeoService {
	return &GeoService{
		geoAdapter: geoAdapter,
	}
}

// DetectAddressInfoRecursive Рекурсивная функция. Если местоположение не найдено, возвращает МСК
func (gs GeoService) DetectAddressInfoRecursive(lat float64, lon float64) (*dto.Location, error) {
	var kladrId string
	var fiasId string
	var typeName string
	geoLocation, err := gs.geoAdapter.DetectAddressInfo(lat, lon)
	// Рекурсия
	if err != nil { // Не смогли определить местоположение. Возвращает дефолтное значение
		if lat == defaultLat && lon == defaultLon { // Вызов на дефолтном местоположении. Больше ничего сделать не можем - ошибка
			return nil, err
		} else { // Пробуем получить дефолтное местоположение
			return gs.DetectAddressInfoRecursive(defaultLat, defaultLon)
		}
	}
	fiasLevel := geoLocation.FiasLevel
	typeCode := fiasLevelMap[fiasLevel]
	if typeCode == "" { // мы не знаем, что это
		return nil, appErrors.InternalError{Err: pkgErrors.WithStack(fmt.Errorf("unkown type of object"))}
	}
	if typeCode != "city" && typeCode != "settlement" {
		if geoLocation.CityFiasId != "" { // это горд
			typeCode = "city"
			fiasId = geoLocation.CityFiasId
			kladrId = geoLocation.CityKladrId
			typeName = geoLocation.CityTypeFull
		} else if geoLocation.SettlementFiasId != "" { // это нас. пункт
			typeCode = "settlement"
			fiasId = geoLocation.SettlementFiasId
			kladrId = geoLocation.SettlementKladrId
			typeName = geoLocation.SettlementTypeFull
		}
	}
	var rusRegionIdStr string
	var rusRegionId int
	if len(geoLocation.TaxOffice) > 0 {
		rusRegionIdStr = geoLocation.TaxOffice[:len(geoLocation.TaxOffice)-2]
		rusRegionId, _ = strconv.Atoi(rusRegionIdStr)
	} else { // Не удалось определить местоположение, невалидный объект
		return gs.DetectAddressInfoRecursive(defaultLat, defaultLon)
	}
	locationResponse := &dto.Location{
		FiasId:         fiasId,
		KladrId:        kladrId,
		CityName:       geoLocation.CityName,
		SettlementName: geoLocation.SettlementName,
		RegionName:     fmt.Sprintf("%s %s", geoLocation.Region, geoLocation.RegionTypeFull),
		AreaName:       fmt.Sprintf("%s %s", geoLocation.Area, geoLocation.AreaTypeFull),
		RegionFiasId:   geoLocation.RegionFiasId,
		RusRegionId:    rusRegionId,
		TypeCode:       typeCode,
		TypeName:       typeName,
	}
	return locationResponse, nil
}

// SuggestAddress Возвращает подсказки по параметрам
func (gs GeoService) SuggestAddress(requestSuggestAddressDto dto.SuggestsRequestDto) (*[]dto.LocationAddress, error) {
	suggests, err := gs.geoAdapter.SuggestAddress(requestSuggestAddressDto)
	if err != nil {
		return nil, err
	}
	var suggestsResponse []dto.LocationAddress
	for _, geoLocation := range suggests {
		locationResponse := buildLocationAddressResponse(&geoLocation)
		suggestsResponse = append(suggestsResponse, locationResponse)
	}
	return &suggestsResponse, nil
}

// FindAddressById Находит адрес по коду КЛАДР или ФИАС
func (gs GeoService) FindAddressById(findAddressByIdRequestDto dto.FindAddressByIdRequestDto) (*dto.LocationAddress, error) {
	locationResponse := dto.LocationAddress{}
	var err error
	geoLocation, err := gs.geoAdapter.FindAddressById(findAddressByIdRequestDto)
	if err != nil {
		return nil, err
	}
	// Строим объект
	locationResponse = buildLocationAddressResponse(geoLocation)
	return &locationResponse, nil
}
