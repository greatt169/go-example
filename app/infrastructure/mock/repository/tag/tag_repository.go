package tag

import (
	uuid "github.com/satori/go.uuid"
	"news-ms/domain/tag/entity"
)

type TagRepositoryMock struct{}

// Получить теги
func (t *TagRepositoryMock) GetTags() []*entity.Tag {
	var tags entity.TagList
	return tags
}

// Удалить все теги
func (t *TagRepositoryMock) RemoveAll() error {
	return nil
}

// Возвращает Id тега по названию
// Если тег не найден, возвращает новый uuid
func (t *TagRepositoryMock) GetTagIdByName(_ string) uuid.UUID {
	return uuid.NewV4()
}
