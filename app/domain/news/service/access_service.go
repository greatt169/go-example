package service

import (
	helpers "github.com/AeroAgency/golang-helpers-lib"
	helpersDto "github.com/AeroAgency/golang-helpers-lib/dto"
	"news-ms/domain/news/entity"
	"time"
)

// Сервис для разграничению доступа согласно ролевой модели
type NewsAccess struct {
	access *helpers.Access
}

// Конструктор
func NewNewsAccess() *NewsAccess {
	return &NewsAccess{
		access: helpers.NewAccess(),
	}
}

// Возвращает признак: может ли пользователь просматривать новости-черновики
func (a NewsAccess) CheckCanShowDraftsNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "show_deactivated")
	return res
}

// Возвращает признак: может ли пользователь просматривать новости с отложенной публикацией
func (a NewsAccess) CheckCanShowDelayedNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "show_delayed")
	return res
}

// Проверка: может ли пользователь просматривать новости-черновики. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanShowDraftNewsOrFail(privilegesDto helpersDto.Privileges, newsEntity *entity.News) error {
	isCanShowDraftsNews := a.CheckCanShowDraftsNews(privilegesDto)
	if newsEntity.Active == false && isCanShowDraftsNews == false {
		return a.access.PermissionDenied()
	}
	return nil
}

// Проверка: может ли пользователь просматривать новости с отложенной публикацией. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanShowDelayedNewsOrFail(privilegesDto helpersDto.Privileges, newsEntity *entity.News) error {
	isCanShowDelayedNews := a.CheckCanShowDelayedNews(privilegesDto)
	if newsEntity.ActiveFrom > time.Now().Unix() && isCanShowDelayedNews == false {
		return a.access.PermissionDenied()
	}
	return nil
}

// Проверка: может ли пользователь создвывть новости. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanCreateNewsOrFail(privilegesDto helpersDto.Privileges) error {
	res := a.access.HasScope(privilegesDto.Entities.News, "create")
	if res != true {
		return a.access.PermissionDenied()
	}
	return nil
}

// Возвращает признак: может ли пользователь редактировать опубликованные новости
func (a NewsAccess) checkCanUpdateActiveNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "update_active")
	return res
}

// Возвращает признак: может ли пользователь редактировать новости с отложенной публикацией
func (a NewsAccess) checkCanUpdateDeactivatedNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "update_deactivated")
	return res
}

// Проверка: может ли пользователь редактировать новости. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanUpdateNewsOrFail(privilegesDto helpersDto.Privileges, newsEntity *entity.News) error {
	isCanUpdateActiveNews := a.checkCanUpdateActiveNews(privilegesDto)
	isCanUpdateDeactivatedNews := a.checkCanUpdateDeactivatedNews(privilegesDto)
	if (newsEntity.Active == false && isCanUpdateDeactivatedNews == false) || (newsEntity.Active == true && isCanUpdateActiveNews == false) {
		return a.access.PermissionDenied()
	}
	return nil
}

// Возвращает признак: может ли пользователь удалять опубликованные новости
func (a NewsAccess) checkCanDeleteActiveNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "delete_active")
	return res
}

// Возвращает признак: может ли пользователь удалять новости с отложенной публикацией
func (a NewsAccess) checkCanDeleteDeactivatedNews(privilegesDto helpersDto.Privileges) bool {
	res := a.access.HasScope(privilegesDto.Entities.News, "delete_deactivated")
	return res
}

// Проверка: может ли пользователь удалять новости. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanDeleteNewsOrFail(privilegesDto helpersDto.Privileges, newsEntity *entity.News) error {
	isCanDeleteActiveNews := a.checkCanDeleteActiveNews(privilegesDto)
	isCanDeleteDeactivatedNews := a.checkCanDeleteDeactivatedNews(privilegesDto)
	if (newsEntity.Active == false && isCanDeleteDeactivatedNews == false) || (newsEntity.Active == true && isCanDeleteActiveNews == false) {
		return a.access.PermissionDenied()
	}
	return nil
}

// Проверка: может ли пользователь фильтровать новости. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanFilterNewsOrFail(privilegesDto helpersDto.Privileges) error {
	res := a.access.HasScope(privilegesDto.Entities.News, "filter")
	if res != true {
		return a.access.PermissionDenied()
	}
	return nil
}

// Проверка: может ли пользователь просматривать детальную страницу новости. В случае запрета возвращается ошибка доступа
func (a NewsAccess) CheckCanShowNewsDetail(privilegesDto helpersDto.Privileges) error {
	res := a.access.HasScope(privilegesDto.Entities.News, "show_news_detail")
	if res != true {
		return a.access.PermissionDenied()
	}
	return nil
}
