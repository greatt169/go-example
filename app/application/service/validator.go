package service

import (
	"fmt"
	"geo/application/dto"
	appErrors "geo/infrastructure/errors"
	pkgErrors "github.com/pkg/errors"
	"github.com/thedevsaddam/govalidator"
)

// Правила валидации входных данных
type Validator struct{}

// Возвращает правила валидации для базовых параметров запроса получения подсказок
func (vs *Validator) getSuggestAddressBaseMap() map[string][]string {
	mapRules := map[string][]string{
		"count":     {"required", "numeric_between:1,100"},
		"fromBound": {"in:,region,city,street,house,settlement"},
		"toBound":   {"in:,region,city,street,house,settlement"},
	}
	return mapRules
}

// Возвращает правила валидации параметров ограничения по родителю для запроса получения подсказок
func (vs *Validator) getSuggestAddressLocationsMap() map[string][]string {
	mapRules := map[string][]string{
		"regionFiasId":     {"regex:(^$|^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$)"},
		"cityFiasId":       {"regex:(^$|^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$)"},
		"settlementFiasId": {"regex:(^$|^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$)"},
		"streetFiasId":     {"regex:(^$|^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$)"},
		"kladrId":          {"regex:^$|^-?[0-9]+$"},
	}
	return mapRules
}

// Возвращает правила валидации параметров запроса получения адреса по id КЛАДР/ФИАС
func (vs *Validator) getFindAddressByIdMap() map[string][]string {
	mapRules := map[string][]string{
		"query": {"required", "regex:(^[+-]?([0-9]*\\.?[0-9]+|[0-9]+\\.?[0-9]*)([eE][+-]?[0-9]+)?$|^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$)"},
	}
	return mapRules
}

// Проверяет на коректность параметры запроса получения подсказок по адресам
func (vs *Validator) ValidateSuggestAddressDto(suggestParams dto.SuggestsRequestDto) error {
	opts := govalidator.Options{
		Data:  &suggestParams,                // request object
		Rules: vs.getSuggestAddressBaseMap(), // rules map
	}
	validationErr := validateStruct(opts)
	if validationErr != nil {
		return validationErr
	}
	// Валидируем поля locations
	if suggestParams.Locations != nil && len(suggestParams.Locations) > 0 {
		for i, location := range suggestParams.Locations {
			elementNum := i + 1
			opts := govalidator.Options{
				Data:  &location,                          // request object
				Rules: vs.getSuggestAddressLocationsMap(), // rules map
			}
			v := govalidator.New(opts)
			validationResult := v.ValidateStruct()
			var firstErrorMsg string
			if len(validationResult) > 0 {
				for _, item := range validationResult {
					firstErrorMsg = item[0]
					break
				}
				return appErrors.BadRequestError{
					Err:     pkgErrors.WithStack(fmt.Errorf(fmt.Sprintf("Element #%v validation error", elementNum))),
					Message: firstErrorMsg,
				}
			}
		}
	}
	return nil
}

// Проверяет на коректность параметры запроса получения адреса по id КЛАДР/ФИАС
func (vs *Validator) ValidateFindAddressByIdDto(findAddressByIdParams dto.FindAddressByIdRequestDto) error {
	opts := govalidator.Options{
		Data:  &findAddressByIdParams,     // request object
		Rules: vs.getFindAddressByIdMap(), // rules map
	}
	validationErr := validateStruct(opts)
	if validationErr != nil {
		return validationErr
	}
	return nil
}

// Общий метод для проверки корректности структуры
func validateStruct(opts govalidator.Options) error {
	v := govalidator.New(opts)
	validationResult := v.ValidateStruct()
	var firstErrorMsg string
	if len(validationResult) > 0 {
		for _, item := range validationResult {
			firstErrorMsg = item[0]
			break
		}
		return appErrors.BadRequestError{
			Err:     pkgErrors.WithStack(fmt.Errorf("validation error")),
			Message: firstErrorMsg,
		}
	}
	return nil
}
