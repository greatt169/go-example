package service

import (
	"fmt"
	"geo/application/dto"
	infrGeoDto "geo/infrastructure/dto"
	"strconv"
)

// Возвращает название типа по коду
func getTypeNameByCode(code string, geoLocation *infrGeoDto.LocationAddressObjectDto) string {
	switch code {
	case "country":
		return ""
	case "region":
		return geoLocation.RegionTypeFull
	case "area":
		return geoLocation.AreaTypeFull
	case "city":
		return geoLocation.CityTypeFull
	case "city_district":
		return ""
	case "settlement":
		return geoLocation.SettlementTypeFull
	case "street":
		return geoLocation.StreetTypeFull
	case "house":
		return geoLocation.HouseTypeFull
	case "flat":
		return geoLocation.FlatTypeFull
	case "plan_structure":
		return ""
	case "unknown":
		return ""
	}
	return ""
}

// Собираем объект адреса на осоновании локации, которая вернулась из адаптера dadata
func buildLocationAddressResponse(geoLocation *infrGeoDto.LocationAddressObjectDto) dto.LocationAddress {
	var rusRegionId int
	fiasLevel := geoLocation.FiasLevel
	typeCode := fiasLevelMap[fiasLevel]
	taxOffice := geoLocation.TaxOffice
	if len(taxOffice) > 0 {
		rusRegionIdStr := taxOffice[:len(taxOffice)-2]
		rusRegionId, _ = strconv.Atoi(rusRegionIdStr)
	} else {
		rusRegionId = 0
	}
	locationResponse := dto.LocationAddress{
		FiasId:         geoLocation.FiasId,
		KladrId:        geoLocation.KladrId,
		CityName:       geoLocation.CityName,
		SettlementName: geoLocation.SettlementName,
		RegionName:     fmt.Sprintf("%s %s", geoLocation.Region, geoLocation.RegionTypeFull),
		AreaName:       fmt.Sprintf("%s %s", geoLocation.Area, geoLocation.AreaTypeFull),
		StreetName:     fmt.Sprintf("%s %s", geoLocation.Street, geoLocation.StreetTypeFull),
		HouseName:      fmt.Sprintf("%s %s", geoLocation.House, geoLocation.HouseTypeFull),
		FlatName:       fmt.Sprintf("%s %s", geoLocation.Flat, geoLocation.FlatTypeFull),
		RegionFiasId:   geoLocation.RegionFiasId,
		RusRegionId:    rusRegionId,
		TypeCode:       typeCode,
		TypeName:       getTypeNameByCode(typeCode, geoLocation),
		Value:          geoLocation.Value,
	}
	return locationResponse
}
