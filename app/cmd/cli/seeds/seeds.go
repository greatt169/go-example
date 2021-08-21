package main

import (
	cache "github.com/AeroAgency/golang-bigcache-lib"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"news-ms/domain/news/entity"
	tagEntity "news-ms/domain/tag/entity"
	minioRepository "news-ms/infrastructure/persistence/awss3"
	"news-ms/infrastructure/persistence/postgres"
	"news-ms/infrastructure/registry"
	"time"
)

func main() {

	// Инициализация контейнера зависимостей
	ctn, err := registry.NewContainer()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}

	// Получение объекта для работы с БД из контейнера
	db := ctn.Resolve("db").(*postgres.Database)

	// Запуск автомиграций
	err = db.Automigrate()
	// получаем необходимые репозитории
	newsRepository := ctn.Resolve("news_repository").(*postgres.NewsRepository)
	fileRepository := ctn.Resolve("file_repository").(*postgres.FileRepository)
	fileStorageRepository := ctn.Resolve("file_storage_repository").(*minioRepository.FileStorageRepository)
	tagRepository := ctn.Resolve("tag_repository").(*postgres.TagRepository)

	// Получение кэш-сервиса из контейнера
	cache := ctn.Resolve("cache").(cache.CacheInterface)

	// очистка файлов
	err = fileStorageRepository.RemoveFolder("news")
	if err != nil {
		log.Fatalf("failed to clear files for news: %v", err)
	}

	// Очистка всех данных в бд сервиса
	err = newsRepository.RemoveAll()
	if err != nil {
		log.Fatalf("failed to remove all news: %v", err)
	}
	tagRepository.RemoveAll()
	if err != nil {
		log.Fatalf("failed to remove all tags: %v", err)
	}

	// Получение демо-новостей для заполнения
	news := getNewsList()
	for _, v := range news {

		// сохраняем файлы в хранилище
		files, err := fileStorageRepository.SaveNewsFiles(&v)
		if err != nil {
			log.Fatalf("failed to save files for news: %v", err)
		}

		// сохраняем новость
		strip := bluemonday.StripTagsPolicy()
		v.TextSearch = strip.Sanitize(v.Text)
		err = newsRepository.Create(v)
		if err != nil {
			log.Fatalf("failed to create item of news: %v", err)
		}

		// Сохраняем файлы
		err = fileRepository.SaveNewsFiles(v.Id, files)
		if err != nil {
			log.Fatalf("failed to save files to db: %v", err)
		}

	}

	// Получение демо-акций для заполнения

	// Сброс кэша
	cache.Clear()
}

// Возвращает timestamp для переданной даты в формате "dd.mm.yyyy"
func getDateTs(date string) int64 {
	layout := "02.01.2006"
	t, _ := time.Parse(layout, date)
	return t.Unix()
}

// Получение демо-новостей для заполнения
func getNewsList() []entity.News {
	var newsList = []entity.News{
		{
			Id:          helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
			Name:        "Новость 1",
			Author:      "Пикин А.С.",
			Active:      true,
			ActiveFrom:  time.Now().Unix(),
			DateCreate:  time.Now().Unix(),
			IsImportant: true,
			Text:        "тестовый текст 1",
			TextJson:    "",
			UserId:      "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f38"),
					Name: "тег",
				},
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f39"),
					Name: "хэштег",
				},
			},
			Files: []entity.File{
				{
					Name:   "Суперновость, файл 1",
					Ext:    "txt",
					Base64: "SGVsbG8gV29ybGQK",
				},
				{
					Name:   "Суперновость, файл 2",
					Ext:    "txt",
					Base64: "SGVsbG8gV29ybGQK",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("ff36d753-7aa9-4f28-8fd3-fb89e1aa463d"),
			Name:       "Новость 2",
			Author:     "Петров И.И.",
			Active:     true,
			ActiveFrom: getDateTs("27.06.2021"),
			DateCreate: time.Now().Unix(),
			Text:       "вася тестовый текст 2",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
					Name: "вася",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("0c06682b-771e-432a-b7c3-02b676b7e920"),
			Name:       "Новость уникальная 2989 тест",
			Author:     "Тестерович И.И.",
			Active:     false,
			ActiveFrom: getDateTs("27.06.2021"),
			DateCreate: time.Now().Unix(),
			Text:       "вася вася уникальный текст",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
					Name: "вася",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("55f4c68a-ccdf-4150-b3a2-9ff2e1585851"),
			Name:       "тег Новость уникальная 2006",
			Author:     "Тестерович И.И.",
			Active:     false,
			ActiveFrom: getDateTs("27.06.2021"),
			DateCreate: time.Now().Unix(),
			Text:       "вася вася уникальный текст",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
					Name: "вася",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("fab483c5-e62e-4788-bac3-3ec23b20c94b"),
			Name:       "Новость 3",
			Author:     "Сидоров С.И.",
			Active:     false,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "Черновик",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f42"),
					Name: "петя",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("27db8d9e-af51-4962-b9a4-bdbdd7567a9d"),
			Name:       "Новость 4",
			Author:     "Алексеев И.И.",
			Active:     true,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "тестовый текст 4",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36888",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f38"),
					Name: "тег",
				},
			},
		},
		// для тестов поиска новости
		{
			Id:         helpers.GetUuidByString("6c1720bc-b87b-47d5-b593-d7090a062f37"),
			Name:       "Конечно Вася",
			Author:     "Курочкин И.В.",
			Active:     true,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "ну кто его не знает",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Files: []entity.File{
				{
					Name:   "Васин файл",
					Ext:    "txt",
					Base64: "SGVsbG8gV29ybGQK",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("6c1720bc-b87b-47d5-b503-d7090a062f37"),
			Name:       "Стиляга из Москвы",
			Author:     "Курочкин И.В.",
			Active:     true,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "Вася, ну кто же ещё",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
		},
		{
			Id:         helpers.GetUuidByString("6c1720bc-b87b-47d5-b593-d7090a062f07"),
			Name:       "Можные ботинки",
			Author:     "Курочкин И.В.",
			Active:     true,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "Бла-бла-бла",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b80b-47d5-b593-d7090a062f38"),
					Name: "вася",
				},
			},
		},
		{
			Id:         helpers.GetUuidByString("6c1720bc-b87b-47d5-b593-d7090a062f39"),
			Name:       "Лялялялялялялялялляляляяяяяяяяяяя тополялялялялялялялял со скидкой",
			Author:     "Курочкин И.В.",
			Active:     true,
			ActiveFrom: time.Now().Unix(),
			DateCreate: time.Now().Unix(),
			Text:       "<h2 style=\"text-align:left;\"><span style=\"color: rgb(0,0,0);background-color: rgb(255,255,255);font-size: 24px;font-family: DauphinPlain;\">Что такое Lorem Ipsum?</span></h2>\n<p style=\"text-align:justify;\"><span style=\"color: rgb(0,0,0);background-color: rgb(255,255,255);font-size: 14px;font-family: Open Sans\", Arial, sans-serif;\"><strong>Lorem Ipsum</strong></span> <span style=\"color: rgb(0,0,0);background-color: rgb(255,255,255);font-size: 14px;font-family: Open Sans\", Arial, sans-serif;\">- это текст-\"рыба\", часто используемый в печати и вэб-дизайне. Lorem Ipsum является стандартной \"рыбой\" для текстов на латинице с начала XVI века. В то время некий безымянный печатник создал большую коллекцию размеров и форм шрифтов, используя Lorem Ipsum для распечатки образцов. Lorem Ipsum не только успешно пережил без заметных изменений пять веков, но и перешагнул в электронный дизайн. Его популяризации в новое время послужили публикация листов Letraset с образцами Lorem Ipsum в 60-х годах и, в более недавнее время, программы электронной вёрстки типа Aldus PageMaker, в шаблонах которых используется Lorem Ipsum.</span>&nbsp;</p>\n<img src=\"data:image/png;base64,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\" alt=\"undefined\" style=\"height: auto;width: auto\"/>\n<p></p>\n",
			TextJson:   "",
			UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
			Tags: []tagEntity.Tag{
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f42"),
					Name: "петя",
				},
				{
					Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f39"),
					Name: "хэштег",
				},
			},
		},
	}
	return newsList
}
