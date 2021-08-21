package entity

import (
	uuid "github.com/satori/go.uuid"
)

// Модель тега новости
type Tag struct {
	Id   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string    `json:"name"`
}

// Модель список тегов
type TagList []*Tag

// Установка названия таблицы в БД
func (t Tag) TableName() string {
	return "tag"
}
