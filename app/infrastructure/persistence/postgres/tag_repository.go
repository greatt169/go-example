package postgres

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"news-ms/domain/tag/entity"
)

// Репозиторий для работы с тегами
type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db}
}

// Получить теги
func (t *TagRepository) GetTags() []*entity.Tag {
	var tags entity.TagList
	t.db.Find(&tags)
	return tags
}

// Удалить все теги
func (t *TagRepository) RemoveAll() error {
	t.db.Delete(entity.Tag{})
	return nil
}

// Возвращает Id тега по названию
// Если тег не найден, возвращает новый uuid
func (t *TagRepository) GetTagIdByName(tag string) uuid.UUID {
	var tagEntity entity.Tag
	db := t.db.Where("name=?", tag).
		Find(&tagEntity)
	if db.Error != nil {
		return uuid.NewV4()
	}
	return tagEntity.Id
}
