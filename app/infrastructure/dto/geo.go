package dto

type LocationObjects struct {
	// Список объектов
	Objects *[]LocationObjectDto
}

// DTO объекта местоположения (город, населенный пункт)
type LocationObjectDto struct {
	// Код ФИАС
	FiasId string
	// Код Кладр
	KladrId string
	// Код региона РФ
	RusRegionId int
	// Номер ФИАС региона
	RegionFiasId string
	// Название области (например: Тверская)
	Region string
	// Тип области (например: область)
	RegionTypeFull string
	// Название района (например: Некрасовский)
	Area string
	// Название типа района (например: район)
	AreaTypeFull string
	// Название города (например: Ярославль)
	CityName string
	// ФИАС код города
	CityFiasId string
	// КЛАДР код города
	CityKladrId string
	// Название типа города (например: город)
	CityTypeFull string
	// Название населенного пункта (например: Рождество)
	SettlementName string
	// ФИАС код населенного пункта
	SettlementFiasId string
	// КЛАДР код населенного пункта
	SettlementKladrId string
	// Название типа населенного пункта (например: село)
	SettlementTypeFull string
	// Уровень объекта в ФИАС. Определяет тип
	FiasLevel string
	// Код ИФНС для физических лиц
	TaxOffice string
	// 	Адрес одной строкой (как показывается в списке подсказок)
}

// DTO объекта адреса местоположения
type LocationAddressObjectDto struct {
	// Код ФИАС
	FiasId string
	// Код Кладр
	KladrId string
	// Код региона РФ
	RusRegionId int
	// Номер ФИАС региона
	RegionFiasId string
	// Название области (например: Тверская)
	Region string
	// Тип области (например: область)
	RegionTypeFull string
	// Название района (например: Некрасовский)
	Area string
	// Название типа района (например: район)
	AreaTypeFull string
	// Название города (например: Ярославль)
	CityName string
	// ФИАС код города
	CityFiasId string
	// КЛАДР код города
	CityKladrId string
	// Название типа города (например: город)
	CityTypeFull string
	// Название населенного пункта (например: Рождество)
	SettlementName string
	// ФИАС код населенного пункта
	SettlementFiasId string
	// КЛАДР код населенного пункта
	SettlementKladrId string
	// Название типа населенного пункта (например: село)
	SettlementTypeFull string
	// Номер ФИАС улицы
	StreetFiasId string
	// Название улицы
	Street string
	// Тип улицы
	StreetTypeFull string
	// Номер ФИАС дома
	HouseFiasId string
	// Название дома
	House string
	// Тип дома
	HouseTypeFull string
	// Название кваритры
	Flat string
	// Тип кваритры
	FlatTypeFull string
	// Уровень объекта в ФИАС. Определяет тип
	FiasLevel string
	// Код ИФНС для физических лиц
	TaxOffice string
	// 	Адрес одной строкой (как показывается в списке подсказок)
	Value string
}
