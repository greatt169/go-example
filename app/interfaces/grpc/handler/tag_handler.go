package handler

import (
	"context"
	cache "github.com/AeroAgency/golang-bigcache-lib"
	"github.com/AeroAgency/golang-helpers-lib"
	"google.golang.org/grpc/codes"
	tagrepository "news-ms/domain/tag/repository"
	content_v1 "news-ms/interfaces/grpc/proto/v1/news"
)

type TagHandler struct {
	repository     tagrepository.TagRepository
	serializer     *helpers.Serializer
	cache          cache.CacheInterface
	hash           *helpers.Hash
	privileges     *helpers.Privileges
	errorFormatter *helpers.ErrorFormatter
}

// Конструктор
func NewTagHandler(repo tagrepository.TagRepository, cache cache.CacheInterface) *TagHandler {
	return &TagHandler{
		repository:     repo,
		serializer:     &helpers.Serializer{},
		cache:          cache,
		errorFormatter: &helpers.ErrorFormatter{},
		privileges:     helpers.NewPrivileges(),
	}
}

// Получение тега
func (t *TagHandler) Get(ctx context.Context, request *content_v1.EmptyRequest) (*content_v1.TagList, error) {
	// Получение привелегий
	_, err := t.privileges.GetPrivilegesForAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	object := content_v1.TagList{}
	tagList := t.repository.GetTags()
	err = t.serializer.ConvertStructToProto(tagList, &object.Tag)
	if err != nil {
		err = t.errorFormatter.Wrap(err, "Caught error while converting to proto")
		err = t.errorFormatter.ReturnError(codes.Internal, err, "При получении списка акций произошла ошибка.")
		return nil, err
	}
	return &object, nil
}
