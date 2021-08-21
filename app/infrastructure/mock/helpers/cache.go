package helpers

import (
	"errors"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"news-ms/domain/news/entity"
	promoEntity "news-ms/domain/promo/entity"
	"reflect"
)

type CacheMock struct {
	cache map[string]interface{}
}

func (c *CacheMock) Get(key string, object interface{}) interface{} {
	res := c.cache[key]
	return res
}

func (c *CacheMock) Set(key string, object interface{}) error {
	if c.cache == nil {
		c.cache = make(map[string]interface{})
	}
	c.cache[key] = object

	newsList, ok := object.(*entity.NewsList)
	if ok {
		if newsList.Total == 13 {
			return errors.New("Save error")
		}
	}

	promoList, ok := object.(*promoEntity.PromoList)
	if ok {
		expected := &promoEntity.PromoList{
			promoEntity.Promo{
				Id:       helpers.GetUuidByString("bf037366-07c2-44ca-a0a3-15598b282b31"),
				Name:     "Акция 4",
				Author:   "Васильев И.И.",
				Active:   true,
				Text:     "тестовый текст 4",
				TextJson: "{\"blocks\":[{\"key\":\"22o5j\",\"text\":\"Новость векааааааа!\",\"type\":\"unstyled\",\"depth\":0,\"inlineStyleRanges\":[{\"offset\":0,\"length\":13,\"style\":\"BOLD\"}],\"entityRanges\":[],\"data\":{}}],\"entityMap\":{}}",
				Period:   "Период действия акции не ограничен",
			},
		}
		if reflect.DeepEqual(promoList, expected) {
			return errors.New("Save error")
		}
	}
	return nil
}

func (c *CacheMock) Delete(key ...string) error {
	for _, v := range key {
		if v == "error" {
			return errors.New("Delete error")
		}
		delete(c.cache, v)
	}
	return nil
}

func (c *CacheMock) Clear() error {
	res := c.cache["error"]
	if res == "error" {
		return errors.New("Cache Clear Error")
	}
	c.cache = map[string]interface{}{}
	return nil
}
