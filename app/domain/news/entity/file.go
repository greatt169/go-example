package entity

import (
	uuid "github.com/satori/go.uuid"
)

// Модель файла для отображения
type File struct {
	ID         string `json:"id"`
	UserId     string `json:"userId"`
	DateCreate int64  `json:"dateCreate"`
	Name       string `json:"name"`
	Ext        string `json:"ext"`
	Path       string
	Type       string
	Base64     string    `json:"base64"`
	EntityId   uuid.UUID // Связь с сущностью, к которой привязан файл
}

// Установка имени таблицы в БД
func (f File) TableName() string {
	return "files"
}
