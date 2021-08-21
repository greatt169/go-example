package news

import helpersDto "github.com/AeroAgency/golang-helpers-lib/dto"

// DTO параметров запроса получения списка новостей
type ListRequestDto struct {
	Offset     int64              `json:"offset"`
	Limit      int64              `json:"limit"`
	Sort       string             `json:"sort"`
	Order      string             `json:"order"`
	Filter     *ListRequestFilter `json:"filter,omitempty"`
	Query      string             `json:"query"`
	Tags       []string
	Privileges helpersDto.Privileges
}

type ListRequestFilter struct {

	// Фильтр по активности
	//   1. не передано - опубликованные и черновики
	//   2. active - только опубликованные
	//   3. inactive - только черновики
	Mode string `json:"mode"`

	// Фильтр по пользователю
	UserId string `json:"userId,omitempty"`

	// Фильтр по активности с определенного времени (timestamp).
	// Если передано, то фильтруем новости (Element ActiveFrom < переданного значения)
	ActiveFrom int64 `json:"activeFrom"`

	// Делалась ли рассылка на почту
	// 0. Да
	// 1. Нет
	IsMailed string
}
