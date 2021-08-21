package handler

import (
	"context"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"google.golang.org/grpc/codes"
	dto "news-ms/application/dto/news"
	"news-ms/application/service"
	"news-ms/application/usecases/news"
	"news-ms/domain/news/entity"
	content_v1 "news-ms/interfaces/grpc/proto/v1/news"
)

type NewsHandler struct {
	newsService    *news.NewsService
	serializer     *helpers.Serializer
	errorFormatter *helpers.ErrorFormatter
	validator      *helpers.Validator
	privileges     *helpers.Privileges
	meta           *helpers.Meta
	validatorRules *service.ValidatorRulesMap
	user           *service.User
}

// Конструктор
func NewNewsHandler(newsService *news.NewsService) *NewsHandler {
	return &NewsHandler{
		newsService:    newsService,
		serializer:     &helpers.Serializer{},
		errorFormatter: &helpers.ErrorFormatter{},
		validator:      &helpers.Validator{},
		privileges:     helpers.NewPrivileges(),
		meta:           &helpers.Meta{},
		validatorRules: &service.ValidatorRulesMap{},
		user:           service.NewUser(),
	}
}

// Возвращает список новостей
func (n *NewsHandler) GetNews(ctx context.Context, requestParams *content_v1.NewsRequestParams) (*content_v1.NewsList, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	// Получение DTO из запроса
	var listRequestParamsDto dto.ListRequestDto
	err = n.serializer.ConvertProtoToStruct(requestParams, &listRequestParamsDto)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while converting proto to struct")
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	// Валидация DTO
	err = n.validator.ValidateProto(
		&listRequestParamsDto,
		n.validatorRules.GetGetNewsMap(),
	)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while validating struct")
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	// Присваиваем привелегии
	listRequestParamsDto.Privileges = privileges
	// Получаем данные
	newsList, err := n.newsService.GetNews(listRequestParamsDto)
	if err != nil {
		return nil, err
	}
	// Формируем ответ
	result := content_v1.NewsList{}
	err = n.serializer.ConvertStructToProto(newsList, &result)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while converting to proto")
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	return &result, nil
}

// Получение детальной информации по новости по ID
func (n *NewsHandler) GetOne(ctx context.Context, obj *content_v1.ObjectId) (*content_v1.NewsObject, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	var protoNewsObject content_v1.NewsObject
	// Валидация параметров
	err = n.validator.ValidateProto(
		&map[string]string{
			"id": obj.Id,
		},
		n.validatorRules.GetGetOneMap(),
	)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	// Получение детальной информации по новости
	newsEntity, err := n.newsService.GetOne(obj.Id, privileges)
	if err != nil {
		return nil, err
	}
	// Преобрзование структуры в Proto
	err = n.serializer.ConvertProtoToStruct(newsEntity, &protoNewsObject)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	return &protoNewsObject, nil
}

// Получение детальной информации по новости по slug
func (n *NewsHandler) GetOneBySlug(ctx context.Context, obj *content_v1.ObjectSlug) (*content_v1.NewsObject, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	var protoNewsObject content_v1.NewsObject
	// Валидация параметров
	err = n.validator.ValidateProto(
		&map[string]string{
			"slug": obj.Slug,
		},
		n.validatorRules.GetGetOneBySlugMap(),
	)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	// Получение детальной информации по новости
	newsEntity, err := n.newsService.GetOneBySlug(obj.Slug, privileges)
	if err != nil {
		return nil, err
	}
	// Преобрзование структуры в Proto
	err = n.serializer.ConvertProtoToStruct(newsEntity, &protoNewsObject)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	return &protoNewsObject, nil
}

// Создание новости
func (n *NewsHandler) Create(ctx context.Context, obj *content_v1.RequestNewsObject) (*content_v1.BaseResponse, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	// Сериализация запроса в entity
	var newsEntity entity.News
	err = n.serializer.ConvertProtoToStruct(obj, &newsEntity)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	// Валидация новости
	err = n.validator.ValidateProto(
		&newsEntity,
		n.validatorRules.GetCreateNewsMap(),
	)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	err = n.validateNewsTags(newsEntity)
	if err != nil {
		return nil, err
	}
	// Получаем автора из токена
	newsEntity.Author = n.user.GetNameByToken(n.meta.GetParam(ctx, "access-token"))
	// Сохранение новости
	err = n.newsService.Create(newsEntity, privileges)
	if err != nil {
		return nil, err
	}
	return &content_v1.BaseResponse{Success: true}, nil
}

// Обновление новости
func (n *NewsHandler) Update(ctx context.Context, obj *content_v1.RequestNewsObject) (*content_v1.BaseResponse, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	// Сериализация запроса в entity
	var newsEntity entity.News
	err = n.serializer.ConvertProtoToStruct(obj, &newsEntity)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.Internal, err, "")
		return nil, err
	}
	// Валидация
	err = n.validator.ValidateProto(
		&newsEntity,
		n.validatorRules.GetUpdateNewsMap(),
	)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	err = n.validateNewsTags(newsEntity)
	if err != nil {
		return nil, err
	}
	// Обновление новости
	err = n.newsService.Update(newsEntity, privileges)
	if err != nil {
		return nil, err
	}
	return &content_v1.BaseResponse{Success: true}, nil
}

// Удаление новости
func (n *NewsHandler) Delete(ctx context.Context, obj *content_v1.ObjectId) (*content_v1.BaseResponse, error) {
	// Получение привелегий
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	// Валидация
	err = n.validator.ValidateProto(
		&map[string]string{
			"id": obj.Id,
		},
		n.validatorRules.GetDeleteMap(),
	)
	if err != nil {
		err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
		return nil, err
	}
	// Удаление новости
	err = n.newsService.Delete(obj.Id, privileges)
	if err != nil {
		return nil, err
	}
	return &content_v1.BaseResponse{Success: true}, nil
}

// Получение ссылки на скачивание вложения новости
func (n *NewsHandler) GetFileLink(ctx context.Context, fileId *content_v1.FileId) (*content_v1.FileLink, error) {
	//Получаем привилегии
	privileges, err := n.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while getting privileges from meta")
		return nil, n.errorFormatter.ReturnError(codes.PermissionDenied, err, "")
	}
	//Валидируем входные данные
	validateRules := n.validatorRules.GetGetOneMap()
	err = n.validator.ValidateProto(
		&map[string]string{
			"id": fileId.Id,
		},
		validateRules,
	)
	if err != nil {
		err = n.errorFormatter.Wrap(err, "Caught error while validating request")
		return nil, n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
	}

	fileInfo, err := n.newsService.GetAttach(fileId.Id, privileges)
	if err != nil {
		return nil, err
	}
	return &content_v1.FileLink{Link: fileInfo.Link}, nil
}

// Валидация тегов
func (n *NewsHandler) validateNewsTags(newsEntity entity.News) error {
	for _, tag := range newsEntity.Tags {
		err := n.validator.ValidateProto(
			&tag,
			n.validatorRules.GetCreateNewsTagMap(),
		)
		if err != nil {
			err = n.errorFormatter.ReturnError(codes.FailedPrecondition, err, "")
			return err
		}
	}
	return nil
}
