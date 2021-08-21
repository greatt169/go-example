package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/machiel/slugify"
	uuid "github.com/satori/go.uuid"
	"news-ms/application/dto/news"
	"news-ms/domain/news/entity"
	"time"
)

type NewsRepository struct {
	db           *gorm.DB
	defaultLimit int64
	defaultSort  string
	defaultOrder string
}

// Репозиторий для работы с новосятми
func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		db,
		20,
		"active_from",
		"desc",
	}
}

// Получить список новостей
func (n *NewsRepository) GetNews(dto news.ListRequestDto) *entity.NewsList {
	// Список новостей
	var news entity.NewsList
	var db *gorm.DB
	db = n.db
	limit := dto.Limit
	offset := dto.Offset
	if limit == 0 {
		limit = n.defaultLimit
	}
	query := dto.Query
	tags := dto.Tags
	// Join тегов новостей
	if len(tags) > 0 || len([]rune(query)) > 0 {
		db = db.Joins("left outer join news_tags on news.id = news_tags.news_id")
		db = db.Joins("left outer join tag on tag.id = news_tags.tag_id")
	}
	// Поиск по тегам
	if len(tags) > 0 {
		db = n.addParamsForSearchByTag(db, tags)
	}

	// Уникальность новостей
	db = db.Select("DISTINCT news.*")

	// Поиск по подстроке
	if query != "" {
		db = n.addParamsForSearchByQuery(db, query)
		if dto.Sort == "" {
			dto.Sort = n.getParamsForSortByRank(query)
			db = db.Select("DISTINCT on (" + n.getParamsForSortByRank(query) + ",news.id) news.*")
		} else {
			db = db.Select("DISTINCT on (" + dto.Sort + ",news.id) news.*")
		}
		if dto.Order == "" {
			dto.Order = "desc"
		}
	}
	// Простановка сортировки
	if dto.Order == "" {
		dto.Order = n.defaultOrder
	}
	if dto.Sort == "" {
		dto.Sort = n.defaultSort
	}
	db = db.Order(dto.Sort + " " + dto.Order + ", id desc")
	// Фильтрация
	filter := dto.Filter
	if filter != nil {
		filterActive := filter.Mode == "active"
		filterUserId := filter.UserId
		filterActiveFrom := filter.ActiveFrom
		filterIsMailed := filter.IsMailed
		if filter.Mode != "" {
			db = db.Where("active = ?", filterActive)
		}
		if filterUserId != "" {
			db = db.Where("user_id = ?", filterUserId)
		}
		if filterActiveFrom != 0 {
			db = db.Where("active_from < ?", filterActiveFrom)
		}
		if filterIsMailed != "" {
			db = db.Where("is_mailed = ?", filterIsMailed)
		}
	}
	var count int64
	db = db.Preload("Tags")
	db = db.Preload("FilesInfo")

	// Подсчет total
	db.Limit(1).Find(&news.News).Count(&count)
	// Поиск
	db.Limit(limit).Offset(offset).Find(&news.News)
	news.Total = count
	return &news
}

// Получить новость
func (n *NewsRepository) GetOne(id uuid.UUID) (*entity.News, error) {
	var news entity.News
	db := n.db.Where("id=?", id).Set("gorm:auto_preload", true).Find(&news)
	if db.Error != nil {
		return nil, db.Error
	}
	err := n.db.Model(&news).Association("Tags").Find(&news.Tags).Error
	if err != nil {
		return nil, err
	}
	err = n.db.Model(&news).Association("FilesInfo").Find(&news.FilesInfo).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// Получить новость
func (n *NewsRepository) GetOneBySlug(slug string) (*entity.News, error) {
	var news entity.News
	db := n.db.
		Where("slug=?", slug).
		Where("active = ?", true).
		Where("active_from < ?", time.Now().Unix()).
		Set("gorm:auto_preload", true).Find(&news)
	if db.Error != nil {
		return nil, db.Error
	}
	err := n.db.Model(&news).Association("Tags").Find(&news.Tags).Error
	if err != nil {
		return nil, err
	}
	err = n.db.Model(&news).Association("FilesInfo").Find(&news.FilesInfo).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// Обновить новость
func (n *NewsRepository) Update(news *entity.News) error {
	active := false
	if news.Active == true {
		active = true
	}
	db := n.db.Model(&news).Updates(map[string]interface{}{
		"Id":          news.Id,
		"Name":        news.Name,
		"Active":      active,
		"Text":        news.Text,
		"TextJson":    news.TextJson,
		"ActiveFrom":  news.ActiveFrom,
		"IsImportant": news.IsImportant,
		"IsMailed":    news.IsMailed,
	})
	db.Model(&news).Association("Tags").Replace(news.Tags)
	return db.Error
}

// Создать новость
func (n *NewsRepository) Create(news entity.News) error {
	// generating slug
	if news.Slug == "" {
		n.generateSlug(&news)
	}

	db := n.db.Create(&news)
	return db.Error
}

// генерация уникального слага
func (n NewsRepository) generateSlug(news *entity.News) {
	for i := 0; true; i++ {
		if i == 0 {
			// для первого круга
			news.Slug = slugify.Slugify(news.Name)
		} else {
			// если уже есть такой slug, добавляем к нему цифры
			news.Slug = slugify.Slugify(
				fmt.Sprintf("%s-%d", news.Name, i),
			)
		}
		// проверка на существование слага
		notFound := n.db.Where("slug = ?", news.Slug).
			Not("id", news.Id).
			First(&entity.News{}).
			RecordNotFound()

		if notFound {
			return
		}
	}
}

// Удалить новость
func (n *NewsRepository) Delete(id uuid.UUID) error {
	news := entity.News{Id: id}
	db := n.db.Unscoped().Delete(&news)
	db.Model(&news).Association("Tags").Replace(news.Tags)
	return db.Error
}

// Удалить все новости
func (n *NewsRepository) RemoveAll() error {
	n.db.Delete(entity.News{})
	n.db.Delete(entity.NewsTags{})
	return nil
}

// Добавляет условия для поиска по переданным тегам
func (n *NewsRepository) addParamsForSearchByTag(db *gorm.DB, tags []string) *gorm.DB {
	// ищем новости по тегам
	db = db.Where("tag.name in (?)", tags)
	return db
}

// Добавляет условия для поиска по переданной подстроке
func (n *NewsRepository) addParamsForSearchByQuery(db *gorm.DB, query string) *gorm.DB {
	sqlQuery :=
		"setweight(to_tsvector('russian', news.name), 'A') || setweight(to_tsvector('russian', text_search), 'B') @@ plainto_tsquery('russian', ?) OR news.name ILIKE ? OR news.text_search ILIKE ? OR tag.name in (?)"
	db = db.Where(
		sqlQuery, query,
		"%"+query+"%",
		"%"+query+"%",
		query,
	)
	return db
}

// Добавляет условия для сортировки по релевантности запросу
func (n *NewsRepository) getParamsForSortByRank(query string) string {
	return fmt.Sprintf(
		"ts_rank(setweight(to_tsvector('russian', news.name), 'A') || setweight(to_tsvector('russian', text), 'B'), plainto_tsquery('russian', '%s'))",
		query,
	)
}
