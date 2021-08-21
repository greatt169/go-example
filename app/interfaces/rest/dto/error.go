package dto

type Error struct {
	// символьный код ошибки
	Error string `json:"applicationErrorCode"`
	// текст ошибки
	Message string `json:"message"`
	// код ошибки
	Code int `json:"-"`
	// Дополнительное описание ошибки (исходная ошибка)
	Debug string `json:"debug"`
}
