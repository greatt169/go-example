package entity

import (
	uuid "github.com/satori/go.uuid"
	tag "news-ms/domain/tag/entity"
)

const (
	AssociationTags = "Tags"
)

// Модель новости
type News struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key;foreignkey:id" json:"id"`
	Name        string    `json:"title"`
	Author      string    `json:"author"`
	Active      bool      `gorm:"type:boolean;omitempty"`
	DateCreate  int64
	ActiveFrom  int64     `json:"activeFrom"`
	Text        string    `json:"text"`
	TextJson    string    `json:"textJson"`   // Текст для плагина визуального редактора
	TextSearch  string    `json:"textSearch"` // Текст для поиска (без html-тегов)
	UserId      string    `gorm:"type:varchar(128);omitempty"`
	Tags        []tag.Tag `gorm:"many2many:news_tags;"`
	Files       []File    // Массив файлов для загрузки
	FilesInfo   []File    `gorm:"ForeignKey:EntityId;AssociationForeignKey:Id"` // Массив файлов для отображения
	IsImportant bool      `json:"isImportant"`
	IsMailed    bool      `gorm:"type:boolean;omitempty"`
	Slug        string
}

// Модель списка новостей
type NewsList struct {
	News  []News `json:"news,omitempty"`
	Total int64  `json:"total"`
}

// Установка имени таблицы в БД
func (n News) TableName() string {
	return "news"
}
