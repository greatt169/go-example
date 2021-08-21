package dto

// Dto для запроса подсказок
type SuggestsRequestDto struct {
	// Количество возвращаемых объектов
	Count int `json:"count"`
	// Поисковая строка по местоположениям
	Query string `json:"query"`
	// Ограничение по типу местоположения (верхняя граница)
	FromBound string `json:"fromBound"`
	// Ограничение по типу местоположения (нижняя граница). Пример city - только города.
	ToBound string `json:"toBound"`
	// 	Ограничение по родителю (страна, регион, район, город, улица)
	Locations []SuggestsRequestLocationDto `json:"locations"`
}

// Ограничение по родителю (страна, регион, район, город, улица)
type SuggestsRequestLocationDto struct {
	// Ограничение по ФИАС-коду региона родителя
	RegionFiasId string `json:"regionFiasId"`
	// Ограничение по ФИАС-коду города родителя
	CityFiasId string `json:"cityFiasId"`
	// Ограничение по ФИАС-коду населенного пункта родителя
	SettlementFiasId string `json:"settlementFiasId"`
	// Ограничение по ФИАС-коду улицы родителя
	StreetFiasId string `json:"streetFiasId"`
	// Ограничение по КЛАДР-коду родителя
	KladrId string `json:"kladrId"`
	// Ограничение по названию региона родителя
	RegionName string `json:"regionName"`
	// Ограничение по названию города родителя
	CityName string `json:"cityName"`
	// Ограничение по названию населенного пункта родителя
	SettlementName string `json:"settlementName"`
	// Ограничение по названию улицы родителя
	StreetName string `json:"streetName"`
}

// DTO запроса поиска адреса по fias/kladr
type FindAddressByIdRequestDto struct {
	// Id fias/kladr
	Query string `json:"query"`
}
