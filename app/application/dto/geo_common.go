package dto

// Объект местоположения (город, населенный пункт)
type Location struct {
	// Код ФИАС
	FiasId string `json:"fiasId"`
	// Код Кладр
	KladrId string `json:"kladrId"`
	// Код региона РФ
	RusRegionId int `json:"rusRegionId"`
	// Номер ФИАС региона
	RegionFiasId string `json:"regionFiasId"`
	// Название области (например: Тверская область)
	RegionName string `json:"regionName"`
	// Название района (например: Некрасовский район)
	AreaName string `json:"areaName"`
	// Название города (например: Ярославль)
	CityName string `json:"cityName"`
	// Название населенного пункта (например: Рождество)
	SettlementName string `json:"settlementName"`
	// Символьный код типа (Например city)
	TypeCode string `json:"typeCode"`
	// Название типа (Например город)
	TypeName string `json:"typeName"`
}

// Объект адреса местоположения
type LocationAddress struct {
	// Код ФИАС
	FiasId string `json:"fiasId"`
	// Код Кладр
	KladrId string `json:"kladrId"`
	// Код региона РФ
	RusRegionId int `json:"rusRegionId"`
	// Номер ФИАС региона
	RegionFiasId string `json:"regionFiasId"`
	// Название области (например: Тверская область)
	RegionName string `json:"regionName"`
	// Название района (например: Некрасовский район)
	AreaName string `json:"areaName"`
	// Название города (например: Ярославль)
	CityName string `json:"cityName"`
	// Название населенного пункта (например: Рождество)
	SettlementName string `json:"settlementName"`
	// Название Улицы
	StreetName string `json:"streetName"`
	// Название дома
	HouseName string `json:"houseName"`
	// Название квартиры
	FlatName string `json:"flatName"`
	// Символьный код типа (Например city)
	TypeCode string `json:"typeCode"`
	// Название типа (Например город)
	TypeName string `json:"typeName"`
	// Адрес одной строкой (как показывается в списке подсказок)
	Value string `json:"value"`
}
