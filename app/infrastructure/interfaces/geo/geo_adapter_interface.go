package geo

import "geo/infrastructure/dto"
import appDto "geo/application/dto"

// Интерфейс сервиса для работы с  api
type GeoAdapterInterface interface {
	// Возвращает данные о местоположении по координатам
	DetectAddressInfo(lat float64, lon float64) (*dto.LocationObjectDto, error)
	// Возвращает подсказки по араметрам
	SuggestAddress(appDto.SuggestsRequestDto) ([]dto.LocationAddressObjectDto, error)
	// Находит адрес по коду КЛАДР или ФИАС
	FindAddressById(requestDto appDto.FindAddressByIdRequestDto) (*dto.LocationAddressObjectDto, error)
}
