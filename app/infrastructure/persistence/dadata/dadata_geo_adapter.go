package dadata

import (
	"fmt"
	appDto "geo/application/dto"
	"geo/infrastructure/dto"
	appErrors "geo/infrastructure/errors"
	pkgErrors "github.com/pkg/errors"
	"gopkg.in/webdeskltd/dadata.v2"
)

// todo: исправить момент, когда dadata возвращает nil
// Сервис для работы с dadata api
type DadataGeoAdapter struct {
	// Клиент по работе с dadata
	dadataClient *dadata.DaData
}

func (d DadataGeoAdapter) FindAddressById(requestDto appDto.FindAddressByIdRequestDto) (*dto.LocationAddressObjectDto, error) {
	dadataResponse, err := d.dadataClient.AddressByID(requestDto.Query)
	if err != nil {
		return nil, err
	}
	if dadataResponse == nil {
		return nil, appErrors.InternalError{Err: pkgErrors.WithStack(fmt.Errorf("bad answer from geo api service"))}
	}
	locationDto := convertDadataLocationResponseToLocationAddressObjectDto(dadataResponse)
	return &locationDto, nil
}

// Конструктор
func NewDadataGeoAdapter(dadataClient *dadata.DaData) *DadataGeoAdapter {
	return &DadataGeoAdapter{
		dadataClient: dadataClient,
	}
}

// Возвращает данны о местоположении по координатам
func (d DadataGeoAdapter) DetectAddressInfo(lat float64, lon float64) (*dto.LocationObjectDto, error) {
	dadataResponse, err := d.dadataClient.GeolocateAddress(
		dadata.GeolocateRequest{
			Lat:   float32(lat),
			Lon:   float32(lon),
			Count: 1,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(dadataResponse) == 0 {
		return nil, appErrors.InternalError{Err: pkgErrors.WithStack(fmt.Errorf("bad answer from geo api service"))}
	}
	locationObject := dadataResponse[0]
	locationDto := convertDadataLocationResponseToLocationObjectDto(locationObject)
	return locationDto, nil
}

// Возвращает список подсказок по адресам
func (d DadataGeoAdapter) SuggestAddress(requestDto appDto.SuggestsRequestDto) ([]dto.LocationAddressObjectDto, error) {
	query := requestDto.Query
	fromBound := requestDto.FromBound
	toBound := requestDto.ToBound
	var locations []dadata.SuggestRequestParamsLocation
	if requestDto.Locations != nil && len(requestDto.Locations) > 0 {
		for _, requestLocation := range requestDto.Locations {
			l := dadata.SuggestRequestParamsLocation{
				KladrID:          requestLocation.KladrId,
				Region:           requestLocation.RegionName,
				RegionFiasID:     requestLocation.RegionFiasId,
				City:             requestLocation.CityName,
				CityFiasID:       requestLocation.CityFiasId,
				Settlement:       requestLocation.SettlementName,
				SettlementFiasID: requestLocation.SettlementFiasId,
				Street:           requestLocation.SettlementName,
				StreetFiasID:     requestLocation.StreetFiasId,
			}
			locations = append(locations, l)
		}
	}
	dadataResponse, err := d.dadataClient.SuggestAddresses(
		dadata.SuggestRequestParams{
			Query:     query,
			Count:     requestDto.Count,
			Locations: locations,
			FromBound: dadata.SuggestBound{
				Value: dadata.BoundValue(fromBound),
			},
			ToBound: dadata.SuggestBound{
				Value: dadata.BoundValue(toBound),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var suggests []dto.LocationAddressObjectDto
	for _, locationObject := range dadataResponse {
		locationDto := convertDadataLocationResponseToLocationAddressObjectDto(&locationObject)
		suggests = append(suggests, locationDto)
	}
	return suggests, nil
}
