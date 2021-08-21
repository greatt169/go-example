package repository

import (
	uuid "github.com/satori/go.uuid"
	"news-ms/application/dto/news"
	"news-ms/domain/news/entity"
)

// Интерфейс репозитория для работы с новосятми
type NewsRepository interface {

	// NewNewsRepository

	// Получить список новостей
	GetNews(request news.ListRequestDto) *entity.NewsList

	// Получить новость
	GetOne(id uuid.UUID) (*entity.News, error)

	// Получить новость по slug
	GetOneBySlug(slug string) (*entity.News, error)

	// Обновить новость
	Update(news *entity.News) error

	// Создать новость
	Create(news entity.News) error

	// Удалить новость
	Delete(id uuid.UUID) error

	// Удалить все новости
	RemoveAll() error
}
