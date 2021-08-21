package news

import (
	"errors"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	uuid "github.com/satori/go.uuid"
	"news-ms/application/dto/news"
	"news-ms/domain/news/entity"
	tagEntity "news-ms/domain/tag/entity"
	"time"
)

type NewsRepositoryMock struct{}

func (n *NewsRepositoryMock) GetOneBySlug(slug string) (*entity.News, error) {
	if slug == "6c1729bc-b87b-47d5-b593-d7090a062f37" {
		return &entity.News{
			Slug:   slug,
			Active: true,
		}, nil
	}
	if slug == "6c1729bc-b87b-47d5-b593-d7090a062f48" {
		return &entity.News{
			Slug:   slug,
			Active: true,
			Files: []entity.File{
				{
					Name:   "Суперновость, файл 1",
					Ext:    "invalid",
					Base64: "SGVsbG8gV29ybGQK",
				},
			},
		}, nil
	}
	if slug == "6c1729bc-b87b-47d5-b593-d7090a062f41" {
		return &entity.News{
			Slug:   slug,
			Active: false,
		}, nil
	}
	if slug == "6c1729bc-b87b-47d5-b593-d7090a062f45" {
		return &entity.News{
			Slug:       slug,
			Active:     true,
			ActiveFrom: time.Now().Unix() + 1,
		}, nil
	}
	return nil, errors.New("fail4")
}

// Получить список новостей
func (n *NewsRepositoryMock) GetNews(dto news.ListRequestDto) *entity.NewsList {
	// Список новостей
	var news entity.NewsList
	if dto.Query == "success" {
		return &entity.NewsList{
			News: []entity.News{
				{
					Id:         helpers.GetUuidByString("ff36d753-7aa9-4f28-8fd3-fb89e1aa463d"),
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
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
			},
		}
	}

	if dto.Query == "error_set_cache" {
		return &entity.NewsList{
			Total: 13,
		}
	}
	return &news
}

// Получить новость
func (n *NewsRepositoryMock) GetOne(id uuid.UUID) (*entity.News, error) {
	if id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f37" {
		return &entity.News{
			Id:     id,
			Active: true,
		}, nil
	}
	if id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f48" {
		return &entity.News{
			Id:     id,
			Active: true,
			Files: []entity.File{
				{
					Name:   "Суперновость, файл 1",
					Ext:    "invalid",
					Base64: "SGVsbG8gV29ybGQK",
				},
			},
		}, nil
	}
	if id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f41" {
		return &entity.News{
			Id:     id,
			Active: false,
		}, nil
	}
	if id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f45" {
		return &entity.News{
			Id:         id,
			Active:     true,
			ActiveFrom: time.Now().Unix() + 1,
		}, nil
	}
	return nil, errors.New("fail4")
}

// Обновить новость
func (n *NewsRepositoryMock) Update(news *entity.News) error {
	if news.Name == "Fail" {
		return errors.New("NewsRepository Update Error")
	}
	return nil
}

// Создать новость
func (n *NewsRepositoryMock) Create(news entity.News) error {
	if news.Name == "Fail" {
		return errors.New("NewsRepository Create Error")
	}
	return nil
}

// Удалить новость
func (n *NewsRepositoryMock) Delete(id uuid.UUID) error {
	if id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f45" {
		return errors.New("NewsRepository Delete Error")
	}
	return nil
}

// Удалить все новости
func (n *NewsRepositoryMock) RemoveAll() error {
	return nil
}
