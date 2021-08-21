package repository

import (
	uuid "github.com/satori/go.uuid"
	"news-ms/domain/news/entity"
)

// Интерфейс репозитория для работы с файлами
type FileRepository interface {
	// Создать файл
	Create(file *entity.File) error

	// Удалить файл
	Delete(id string) error

	// Сохраняет переданные файлы, не переданные - удаляет
	SaveNewsFiles(entityUuid uuid.UUID, file []entity.File) error

	// Получить файл по ID
	GetByID(id string) (*entity.File, error)

	// Удалить все файлы
	RemoveAll() error
}
