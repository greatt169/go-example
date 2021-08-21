package news

import (
	cache "github.com/AeroAgency/golang-bigcache-lib"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	helpersDto "github.com/AeroAgency/golang-helpers-lib/dto"
	"github.com/machiel/slugify"
	"github.com/microcosm-cc/bluemonday"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	dto "news-ms/application/dto/news"
	"news-ms/domain/news/entity"
	"news-ms/domain/news/repository"
	"news-ms/domain/news/service"
	tagrepository "news-ms/domain/tag/repository"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Сервис новостей
type NewsService struct {
	newsRepository   repository.NewsRepository
	tagRepository    tagrepository.TagRepository
	newsFilesService *NewsFileService
	cache            cache.CacheInterface
	errorFormatter   *helpers.ErrorFormatter
	hash             *helpers.Hash
	access           *service.NewsAccess
}

// Конструктор
func NewNewsService(
	newsRepository repository.NewsRepository,
	tagRepository tagrepository.TagRepository,
	newsFilesService *NewsFileService,
	cache cache.CacheInterface,
) *NewsService {
	return &NewsService{
		newsRepository:   newsRepository,
		tagRepository:    tagRepository,
		newsFilesService: newsFilesService,
		cache:            cache,
		errorFormatter:   &helpers.ErrorFormatter{},
		access:           service.NewNewsAccess(),
		hash:             &helpers.Hash{},
	}
}

// Внутренний метод получения записи по ID
// Возвращает 404, если запись не найдена
func (n *NewsService) getRecord(recordId string) (*entity.News, error) {
	id, _ := uuid.FromString(recordId)
	newsEntity, err := n.newsRepository.GetOne(id)
	if newsEntity == nil && err != nil {
		err = n.errorFormatter.ReturnError(codes.NotFound, err, "")
		return nil, err
	}
	return newsEntity, nil
}

// Внутренний метод получения записи по slug
// Возвращает 404, если запись не найдена
func (n *NewsService) getRecordBySlug(recordSlug string) (*entity.News, error) {
	newsEntity, err := n.newsRepository.GetOneBySlug(recordSlug)
	if newsEntity == nil && err != nil {
		err = n.errorFormatter.ReturnError(codes.NotFound, err, "")
		return nil, err
	}
	return newsEntity, nil
}

// Возвращает список новостей
func (n *NewsService) GetNews(listRequestParamsDto dto.ListRequestDto) (*entity.NewsList, error) {
	// Если пользователь передает параметры фильтрации по автору или статус публикации проверяем, есть ли у него права
	if listRequestParamsDto.Filter != nil {
		if listRequestParamsDto.Filter.UserId != "" || listRequestParamsDto.Filter.Mode != "" {
			// Проверка доступа
			err := n.access.CheckCanFilterNewsOrFail(listRequestParamsDto.Privileges)
			if err != nil {
				return nil, err
			}
		}
	} else {
		listRequestParamsDto.Filter = &dto.ListRequestFilter{}
	}
	// Может ли пользователь смотреть черновики
	// Если нет, добавляем параметр в фильтр
	isCanShowDrafts := n.access.CheckCanShowDraftsNews(listRequestParamsDto.Privileges)
	if isCanShowDrafts != true {
		listRequestParamsDto.Filter.Mode = "active"
	}

	// Может ли пользователь смотреть отложенно опубликованные
	// Если нет, добавляем параметр в фильтр
	isCanShowDelayed := n.access.CheckCanShowDelayedNews(listRequestParamsDto.Privileges)
	if isCanShowDelayed != true {
		listRequestParamsDto.Filter.ActiveFrom = time.Now().Unix()
	}
	// Проверка результата в кэше
	object := entity.NewsList{}
	key := n.hash.GetHashStringByStruct(listRequestParamsDto)
	resultFromCache := n.cache.Get(key, &object)
	if resultFromCache != nil {
		return resultFromCache.(*entity.NewsList), nil
	}
	// В кеше нет результата - получаем данные из БД и сохраняем в кэш
	query, tags := getSearchParams(listRequestParamsDto.Query)
	// Преобразуем поисковую строку в строку и набор тегов
	listRequestParamsDto.Query = query
	listRequestParamsDto.Tags = tags
	newsList := n.newsRepository.GetNews(listRequestParamsDto)
	err := n.cache.Set(key, newsList)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while set cache")
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	return newsList, nil
}

// Возвращает детальную информацию по новости по ID
func (n *NewsService) GetOne(id string, privileges helpersDto.Privileges) (*entity.News, error) {
	var newsEntity *entity.News
	// Получение детальной информации по новости из БД
	newsEntity, err := n.getRecord(id)
	if err != nil {
		return nil, err
	}
	// Проверка доступа на просмотр черновика в соответствии со свойствами объекта
	err = n.access.CheckCanShowDraftNewsOrFail(privileges, newsEntity)
	if err != nil {
		return nil, err
	}
	// Проверка доступа на просмотр отложенно опубликованного элемента в соответствии со свойствами объекта
	err = n.access.CheckCanShowDelayedNewsOrFail(privileges, newsEntity)
	if err != nil {
		return nil, err
	}
	return newsEntity, nil
}

// Возвращает детальную информацию по новости по slug
func (n *NewsService) GetOneBySlug(slug string, privileges helpersDto.Privileges) (*entity.News, error) {
	// Проверка доступа на просмотр отложенно опубликованного элемента в соответствии со свойствами объекта
	err := n.access.CheckCanShowNewsDetail(privileges)
	if err != nil {
		return nil, err
	}
	var newsEntity *entity.News
	// Получение детальной информации по новости из БД
	newsEntity, err = n.getRecordBySlug(slug)
	if err != nil {
		return nil, err
	}
	return newsEntity, nil
}

// Сохраняет новость
func (n *NewsService) Create(newsEntity entity.News, privileges helpersDto.Privileges) error {
	// Генерация полей для добавление в БД
	newsEntity.Id = uuid.NewV4()
	newsEntity.DateCreate = time.Now().Unix()
	newsEntity.Slug = slugify.Slugify(newsEntity.Name)
	// Связь entity с user
	newsEntity.UserId = privileges.UserId
	// Проверка доступа
	err := n.access.CheckCanCreateNewsOrFail(privileges)
	if err != nil {
		return err
	}
	strip := bluemonday.StripTagsPolicy()
	newsEntity.TextSearch = strip.Sanitize(newsEntity.Text)

	// Сохраняем uuid для тегов
	for index := range newsEntity.Tags {
		if newsEntity.Tags[index].Id == uuid.Nil {
			newsEntity.Tags[index].Id = n.tagRepository.GetTagIdByName(newsEntity.Tags[index].Name)
		}
	}
	err = n.newsRepository.Create(newsEntity)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "При создании новости произошла ошибка.")
		return err
	}
	// Сохраняем файлы новости
	err = n.newsFilesService.SaveFiles(newsEntity)
	if err != nil {
		return err
	}
	// Очистка кэша
	err = n.cache.Clear()
	if err != nil {
		return n.errorFormatter.ReturnError(codes.Internal, err, "")
	}
	return nil
}

// Обновляет новость
func (n *NewsService) Update(newsEntity entity.News, privileges helpersDto.Privileges) error {
	// Связь entity с user
	newsEntity.UserId = privileges.UserId

	// Проверяем новость на существование
	existedNewsEntity, err := n.getRecord(newsEntity.Id.String())
	if err != nil {
		return err
	}
	// Проверка доступа
	err = n.access.CheckCanUpdateNewsOrFail(privileges, existedNewsEntity)
	if err != nil {
		return err
	}
	// Сохраняем uuid для тегов
	for index := range newsEntity.Tags {
		if newsEntity.Tags[index].Id == uuid.Nil {
			newsEntity.Tags[index].Id = n.tagRepository.GetTagIdByName(newsEntity.Tags[index].Name)
		}
	}
	// Сохраняем объект в БД
	err = n.newsRepository.Update(&newsEntity)
	if err != nil {
		return n.errorFormatter.ReturnError(codes.Internal, err, "При обновлении новости произошла ошибка.")
	}
	// Сохраняем файлы новости
	err = n.newsFilesService.SaveFiles(newsEntity)
	if err != nil {
		return err
	}
	// Очистка кэша
	err = n.cache.Clear()
	if err != nil {
		return n.errorFormatter.ReturnError(codes.Internal, err, "")
	}
	return nil
}

// Удаляет новость
func (n *NewsService) Delete(id string, privileges helpersDto.Privileges) error {
	// Проверяем новость на существование
	newsEntity, err := n.getRecord(id)
	if err != nil {
		return err
	}
	// Проверка доступа
	err = n.access.CheckCanDeleteNewsOrFail(privileges, newsEntity)
	if err != nil {
		return err
	}
	// Удаляем объект из БД
	err = n.newsRepository.Delete(newsEntity.Id)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
	}
	// Удаляем файл из хранилища
	err = n.newsFilesService.DeleteFiles(newsEntity.Id)
	if err != nil {
		return err
	}
	// Очистка кэша
	err = n.cache.Clear()
	if err != nil {
		return n.errorFormatter.ReturnError(codes.Internal, err, "")
	}
	return nil
}

// Получает ссылку на скачивание
func (n *NewsService) GetAttach(id string, privileges helpersDto.Privileges) (*dto.FileInfoDto, error) {
	var newsEntity *entity.News
	//Получаем информацию о файле
	fileInfo, err := n.newsFilesService.GetFileInfo(id)
	if err != nil {
		return nil, err
	}
	// Получение детальной информации по новости из БД
	newsEntity, err = n.getRecord(fileInfo.Fields.EntityId.String())
	if err != nil {
		return nil, err
	}
	// Проверка доступа на просмотр черновика в соответствии со свойствами объекта
	err = n.access.CheckCanShowDraftNewsOrFail(privileges, newsEntity)
	if err != nil {
		return nil, err
	}
	// Проверка доступа на просмотр отложенно опубликованного элемента в соответствии со свойствами объекта
	err = n.access.CheckCanShowDelayedNewsOrFail(privileges, newsEntity)
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

// Возвращает параметры для поиска (поисковая строка, теги) в обработанном виде
func getSearchParams(dtoQuery string) (query string, tags []string) {
	// удаляем пробелы с начала и конца строки
	fullQuery := strings.TrimSpace(dtoQuery)
	// заменяем несклько пробелов одним
	r := regexp.MustCompile("\\s+")
	fullQuery = r.ReplaceAllString(fullQuery, " ")
	queryTerms := strings.Split(fullQuery, " ")
	// разделяем запрос на подстроку
	if len(queryTerms) > 0 {
		for _, term := range queryTerms {
			if term == "" {
				continue
			}
			termInitial := string([]rune(term)[0])
			if termInitial == "#" { // это хэшетег
				term = trimFirstRune(term)
				if len([]rune(term)) > 2 {
					tags = append(tags, term)
				}
			} else { // это часть поискового запроса
				query += term + " "
			}
		}
		query = strings.TrimSpace(query)
	}
	return query, tags
}

// Удаляет первый символ
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
