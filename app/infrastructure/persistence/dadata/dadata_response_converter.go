package dadata

import (
	"geo/infrastructure/dto"
	"gopkg.in/webdeskltd/dadata.v2"
)

// Преобразование объектов ответа dadata в DTO
type DadataResponseConverter struct{}

// Преобразует местоположение, которое вернулось из dadata в LocationObjectDto
func convertDadataLocationResponseToLocationObjectDto(locationObject dadata.ResponseAddress) *dto.LocationObjectDto {
	locationDto := &dto.LocationObjectDto{
		// Код ФИАС
		FiasId: locationObject.Data.FiasID,
		// Код Кладр
		KladrId: locationObject.Data.KladrID,
		// Номер ФИАС региона
		RegionFiasId: locationObject.Data.RegionFiasID,
		// Название области (например: Тверская)
		Region: locationObject.Data.Region,
		// Тип области (например: область)
		RegionTypeFull: locationObject.Data.RegionTypeFull,
		// Название района (например: Некрасовский)
		Area: locationObject.Data.Area,
		// Название типа района (например: район)
		AreaTypeFull: locationObject.Data.AreaTypeFull,
		// Название города (например: Ярославль)
		CityName: locationObject.Data.City,
		// ФИАС код города
		CityFiasId: locationObject.Data.CityFiasID,
		// КЛАДР код города
		CityKladrId: locationObject.Data.CityKladrID,
		// Название типа города (например: город)
		CityTypeFull: locationObject.Data.CityTypeFull,
		// Название населенного пункта (например: Рождество)
		SettlementName: locationObject.Data.Settlement,
		// КЛАДР код населенного пункта
		SettlementKladrId: locationObject.Data.SettlementKladrID,
		// Название типа населенного пункта (например: село)
		SettlementTypeFull: locationObject.Data.SettlementTypeFull,
		// Уровень объекта в ФИАС. Определяет тип
		FiasLevel: locationObject.Data.FiasLevel,
		// Код ИФНС для физических лиц
		TaxOffice: locationObject.Data.TaxOffice,
	}
	return locationDto
}

// Преобразует местоположение, которое вернулось из dadata в LocationAddressObjectDto
func convertDadataLocationResponseToLocationAddressObjectDto(locationObject *dadata.ResponseAddress) dto.LocationAddressObjectDto {
	locationDto := dto.LocationAddressObjectDto{
		// Код ФИАС
		FiasId: locationObject.Data.FiasID,
		// Код Кладр
		KladrId: locationObject.Data.KladrID,
		// Номер ФИАС региона
		RegionFiasId: locationObject.Data.RegionFiasID,
		// Название области (например: Тверская)
		Region: locationObject.Data.Region,
		// Тип области (например: область)
		RegionTypeFull: locationObject.Data.RegionTypeFull,
		// Название района (например: Некрасовский)
		Area: locationObject.Data.Area,
		// Название типа района (например: район)
		AreaTypeFull: locationObject.Data.AreaTypeFull,
		// Название города (например: Ярославль)
		CityName: locationObject.Data.City,
		// Название типа города (например: город)
		CityTypeFull: locationObject.Data.CityTypeFull,
		// ФИАС код города
		CityFiasId: locationObject.Data.CityFiasID,
		// КЛАДР код города
		CityKladrId: locationObject.Data.CityKladrID,
		// Название населенного пункта (например: Рождество)
		SettlementName: locationObject.Data.Settlement,
		// КЛАДР код населенного пункта
		SettlementKladrId: locationObject.Data.SettlementKladrID,
		// Название типа населенного пункта (например: село)
		SettlementTypeFull: locationObject.Data.SettlementTypeFull,
		// ФИАС код улицы
		StreetFiasId: locationObject.Data.StreetFiasID,
		// Название улицы
		Street: locationObject.Data.Street,
		// Название типа улицы
		StreetTypeFull: locationObject.Data.StreetTypeFull,
		// ФИАС код дома
		HouseFiasId: locationObject.Data.HouseFiasID,
		// Название дома
		House: locationObject.Data.House,
		// Название типа дома
		HouseTypeFull: locationObject.Data.HouseTypeFull,
		// Название квартиры
		Flat: locationObject.Data.Flat,
		// Название типа квартиры
		FlatTypeFull: locationObject.Data.FlatTypeFull,
		// Уровень объекта в ФИАС. Определяет тип
		FiasLevel: locationObject.Data.FiasLevel,
		// Код ИФНС для физических лиц
		TaxOffice: locationObject.Data.TaxOffice,
		// 	Адрес одной строкой (как показывается в списке подсказок)
		Value: locationObject.Value,
	}
	return locationDto
}
