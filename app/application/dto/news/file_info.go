package news

import "news-ms/domain/news/entity"

// DTO с информацией о файле
type FileInfoDto struct {
	Fields *entity.File
	Link   string
}
