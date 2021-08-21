package repository

import (
	"news-ms/domain/news/entity"
	"time"
)

// Интерфейс для работы с файловым хранилищем
type FileStorageInterface interface {
	// Сохранение файлов из данных сущности
	SaveNewsFiles(News *entity.News) ([]entity.File, error)

	// Получение ссылки
	GetLink(fileData *entity.File, second time.Duration) (string, error)

	// Удаление всех файлов внутри директории хранлища
	RemoveFolder(folderName string) error
}
