package repository

import (
	uuid "github.com/satori/go.uuid"
	"news-ms/domain/tag/entity"
)

// Интерфейс репозитория для работы с тегами
type TagRepository interface {
	// Получить теги
	GetTags() []*entity.Tag

	// Удалить все теги
	RemoveAll() error

	// Возвращает Id тега по названию
	// Если тег не найден, возвращает новый uuid
	GetTagIdByName(tag string) uuid.UUID
}
