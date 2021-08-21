package news

import (
	cache "github.com/AeroAgency/golang-bigcache-lib"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"github.com/AeroAgency/golang-helpers-lib/dto"
	"news-ms/application/dto/news"
	newsdto "news-ms/application/dto/news"
	"news-ms/domain/news/entity"
	"news-ms/domain/news/repository"
	"news-ms/domain/news/service"
	tagEntity "news-ms/domain/tag/entity"
	tagrepository "news-ms/domain/tag/repository"
	helpersMock "news-ms/infrastructure/mock/helpers"
	newsMock "news-ms/infrastructure/mock/repository/news"
	tagMock "news-ms/infrastructure/mock/repository/tag"
	"reflect"
	"testing"
	"time"
)

func TestNewNewsService(t *testing.T) {
	fileRepository := &newsMock.FileRepositoryMock{}
	tagRepository := &tagMock.TagRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	cacheMock := &helpersMock.CacheMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)
	NewNewsService(newsRepository, tagRepository, fileService, cacheMock)
}

func TestNewsService_getRecord(t *testing.T) {
	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		errorFormatter   *helpers.ErrorFormatter
	}
	type args struct {
		recordId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.News
		wantErr bool
	}{
		{
			name: "Case [record returns]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
			},
			args: args{
				recordId: "6c1729bc-b87b-47d5-b593-d7090a062f37",
			},
			want: &entity.News{
				Id:     helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
				Active: true,
			},
			wantErr: false,
		},
		{
			name: "Case [record does not return]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				recordId: "6c1729bc-b87b-47d5-b593-d7090a062ftt",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				errorFormatter:   tt.fields.errorFormatter,
			}
			got, err := n.getRecord(tt.args.recordId)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsService_GetOne(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		access           *service.NewsAccess
		errorFormatter   *helpers.ErrorFormatter
	}
	type args struct {
		id         string
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.News
		wantErr bool
	}{
		{
			name: "Case [success answer]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
			},
			args: args{
				id: "6c1729bc-b87b-47d5-b593-d7090a062f37",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
			},
			want: &entity.News{
				Id:     helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
				Active: true,
			},
			wantErr: false,
		},
		{
			name: "Case [getRecord returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				id: "6c1729bc-b87b-47d5-b593-d7090a062f40",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [access service returns err CheckCanShowDraftNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				id: "6c1729bc-b87b-47d5-b593-d7090a062f41",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [access service returns err CheckCanShowDelayedNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				id: "6c1729bc-b87b-47d5-b593-d7090a062f45",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				access:           tt.fields.access,
				errorFormatter:   tt.fields.errorFormatter,
			}
			got, err := n.GetOne(tt.args.id, tt.args.privileges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsService_GetNews(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)
	cacheMock := &helpersMock.CacheMock{}

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		cache            cache.CacheInterface
		errorFormatter   *helpers.ErrorFormatter
		access           *service.NewsAccess
	}
	type args struct {
		listRequestParamsDto news.ListRequestDto
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.NewsList
		wantErr bool
	}{
		{
			name: "Case [success answer]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				cache:            cacheMock,
			},
			args: args{
				listRequestParamsDto: news.ListRequestDto{
					Query: "success",
					Privileges: dto.Privileges{
						Entities: dto.Entities{
							News: []string{"show_active", "filter"},
						},
					},
				},
			},
			want: &entity.NewsList{
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
			},
			wantErr: false,
		},
		{
			name: "Case [access service returns error CheckCanFilterNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				cache:            cacheMock,
			},
			args: args{
				listRequestParamsDto: news.ListRequestDto{
					Filter: &news.ListRequestFilter{
						UserId: "f22b237d-c92e-49cd-a178-d29334ae9d68",
					},
					Privileges: dto.Privileges{
						Entities: dto.Entities{
							News: []string{"show_active"},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [success result from cache]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				cache: func() *helpersMock.CacheMock {
					cache := &helpersMock.CacheMock{}
					hash := &helpers.Hash{}
					dto := news.ListRequestDto{
						Query:  "from_cache",
						Filter: &news.ListRequestFilter{},
						Privileges: dto.Privileges{
							Entities: dto.Entities{
								News: []string{"show_active", "show_deactivated", "show_delayed", "filter"},
							},
						},
					}
					key := hash.GetHashStringByStruct(dto)
					data := &entity.NewsList{
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
						},
					}
					cache.Set(key, data)
					return cache
				}(),
			},
			args: args{
				listRequestParamsDto: news.ListRequestDto{
					Query: "from_cache",
					Privileges: dto.Privileges{
						Entities: dto.Entities{
							News: []string{"show_active", "show_deactivated", "show_delayed", "filter"},
						},
					},
				},
			},
			want: &entity.NewsList{
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
				},
			},
			wantErr: false,
		},
		{
			name: "Case [error while setting cache]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				listRequestParamsDto: news.ListRequestDto{
					Query: "error_set_cache",
					Privileges: dto.Privileges{
						Entities: dto.Entities{
							News: []string{"show_active", "show_deactivated", "show_delayed", "filter"},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				cache:            tt.fields.cache,
				errorFormatter:   tt.fields.errorFormatter,
				access:           tt.fields.access,
			}
			got, err := n.GetNews(tt.args.listRequestParamsDto)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimFirstRune(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "First char [Case: #qwerty]",
			args: args{
				s: "#qwerty",
			},
			want: "qwerty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimFirstRune(tt.args.s); got != tt.want {
				t.Errorf("trimFirstRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsService_Create(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	tagRepository := &tagMock.TagRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)
	cacheMock := &helpersMock.CacheMock{}

	type fields struct {
		newsRepository   repository.NewsRepository
		tagRepository    tagrepository.TagRepository
		newsFilesService *NewsFileService
		errorFormatter   *helpers.ErrorFormatter
		access           *service.NewsAccess
		cache            cache.CacheInterface
	}
	type args struct {
		newsEntity entity.News
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case [success operation]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"create"},
					},
				},
				newsEntity: entity.News{
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
					DateCreate: time.Now().Unix(),
					Text:       "вася тестовый текст 2",
					TextJson:   "",
					UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
					Files: []entity.File{
						{
							Name:   "success",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
					Tags: []tagEntity.Tag{
						{
							Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
							Name: "вася",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Case [access control returns error CheckCanCreateNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: entity.News{
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
			},
			wantErr: true,
		},
		{
			name: "Case [news repo returns create err]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"create"},
					},
				},
				newsEntity: entity.News{
					Name:       "Fail",
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
			},
			wantErr: true,
		},
		{
			name: "Case [news file saved with err]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"create"},
					},
				},
				newsEntity: entity.News{
					Name:       "Новость века",
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
					Files: []entity.File{
						{
							Name:   "Суперновость, файл 1",
							Ext:    "invalid",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Case [Tag without ID]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"create"},
					},
				},
				newsEntity: entity.News{
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
					DateCreate: time.Now().Unix(),
					Text:       "вася тестовый текст 2",
					TextJson:   "",
					UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
					Files: []entity.File{
						{
							Name:   "success",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
					Tags: []tagEntity.Tag{
						{
							Name: "вася",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Case [Clear cache err after saving]",
			fields: fields{
				newsRepository:   newsRepository,
				tagRepository:    tagRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache: func() *helpersMock.CacheMock {
					cache := &helpersMock.CacheMock{}
					key := "error"
					data := "error"
					cache.Set(key, data)
					return cache
				}(),
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"create"},
					},
				},
				newsEntity: entity.News{
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
					DateCreate: time.Now().Unix(),
					Text:       "вася тестовый текст 2",
					TextJson:   "",
					UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
					Files: []entity.File{
						{
							Name:   "success",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
					Tags: []tagEntity.Tag{
						{
							Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
							Name: "вася",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				tagRepository:    tt.fields.tagRepository,
				errorFormatter:   tt.fields.errorFormatter,
				access:           tt.fields.access,
				cache:            tt.fields.cache,
			}
			if err := n.Create(tt.args.newsEntity, tt.args.privileges); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsService_Update(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)
	cacheMock := &helpersMock.CacheMock{}

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		cache            cache.CacheInterface
		errorFormatter   *helpers.ErrorFormatter
		hash             *helpers.Hash
		access           *service.NewsAccess
	}
	type args struct {
		newsEntity entity.News
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case [success operation]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
					DateCreate: time.Now().Unix(),
					Text:       "вася тестовый текст 2",
					TextJson:   "",
					UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
					Files: []entity.File{
						{
							Name:   "success",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
					Tags: []tagEntity.Tag{
						{
							Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
							Name: "вася",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Case [getRecord returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
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
			},
			wantErr: true,
		},
		{
			name: "Case [access control returns error CheckCanUpdateNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
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
			},
			wantErr: true,
		},
		{
			name: "Case [news repo update returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
					Name:       "Fail",
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
			},
			wantErr: true,
		},
		{
			name: "Case [news file saved with err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
					Name:       "Новость века",
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
					Files: []entity.File{
						{
							Name:   "Суперновость, файл 1",
							Ext:    "invalid",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Case [Clear cache err after saving]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache: func() *helpersMock.CacheMock {
					cache := &helpersMock.CacheMock{}
					key := "error"
					data := "error"
					cache.Set(key, data)
					return cache
				}(),
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: entity.News{
					Id:         helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
					Name:       "Новость 2",
					Author:     "Петров И.И.",
					Active:     true,
					DateCreate: time.Now().Unix(),
					Text:       "вася тестовый текст 2",
					TextJson:   "",
					UserId:     "dc79541b-a854-4c9d-a42f-ec09e0e36887",
					Files: []entity.File{
						{
							Name:   "success",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
					Tags: []tagEntity.Tag{
						{
							Id:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f40"),
							Name: "вася",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				cache:            tt.fields.cache,
				errorFormatter:   tt.fields.errorFormatter,
				hash:             tt.fields.hash,
				access:           tt.fields.access,
			}
			if err := n.Update(tt.args.newsEntity, tt.args.privileges); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsService_Delete(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)
	cacheMock := &helpersMock.CacheMock{}

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		cache            cache.CacheInterface
		errorFormatter   *helpers.ErrorFormatter
		hash             *helpers.Hash
		access           *service.NewsAccess
	}
	type args struct {
		id         string
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case [success operation]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f37",
			},
			wantErr: false,
		},
		{
			name: "Case [getRecord returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f40",
			},
			wantErr: true,
		},
		{
			name: "Case [access control returns error CheckCanDeleteNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f37",
			},
			wantErr: true,
		},
		{
			name: "Case [news repo delete returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f45",
			},
			wantErr: true,
		},
		{
			name: "Case [news file saved with err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache:            cacheMock,
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f48",
			},
			wantErr: true,
		},
		{
			name: "Case [Clear cache err after saving]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
				cache: func() *helpersMock.CacheMock {
					cache := &helpersMock.CacheMock{}
					key := "error"
					data := "error"
					cache.Set(key, data)
					return cache
				}(),
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "6c1729bc-b87b-47d5-b593-d7090a062f37",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				cache:            tt.fields.cache,
				errorFormatter:   tt.fields.errorFormatter,
				hash:             tt.fields.hash,
				access:           tt.fields.access,
			}
			if err := n.Delete(tt.args.id, tt.args.privileges); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsService_GetAttach(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)

	type fields struct {
		newsRepository   repository.NewsRepository
		newsFilesService *NewsFileService
		errorFormatter   *helpers.ErrorFormatter
		access           *service.NewsAccess
	}
	type args struct {
		id         string
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *newsdto.FileInfoDto
		wantErr bool
	}{
		{
			name: "Case [success operation]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "c32267c9-73b8-4ce6-a489-9e2d96d2ccab",
			},
			want: &news.FileInfoDto{
				Fields: &entity.File{
					ID:         "c32267c9-73b8-4ce6-a489-9e2d96d2ccab",
					Name:       "Суперновость, файл 1",
					Ext:        "txt",
					Path:       "news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1",
					UserId:     "dc72578e-da65-4de2-bbc2-8d3fc82c7ebb",
					DateCreate: 1597771527,
					EntityId:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
					Type:       "file",
					Base64:     "SGVsbG8gV29ybGQK",
				},
				Link: "/content-files/news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=123qwe123qwe123qwe%2F20200907%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200907T125915Z&X-Amz-Expires=60&X-Amz-SignedHeaders=host&response-Content-Disposition=attachment%3B%20filename%3D%22%D0%A1%D1%83%D0%BF%D0%B5%D1%80%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D1%8C%2C%20%D1%84%D0%B0%D0%B9%D0%BB%201.txt%22&response-Content-Type=text%2Fplain%3B%20charset%3Dutf-8&X-Amz-Signature=436b8d7973c13bad5467a64c682d3c7cd8c66e758d176b520f2a3278d49888db",
			},
			wantErr: false,
		},
		{
			name: "Case [getFileInfo returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "c32267c9-73b8-4ce6-a489-9e2d96d2ccag",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [getRecord returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "p32267c9-73b8-4ce6-a489-9e2d96d2ccam",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [access control returns err: CheckCanShowDraftNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "k32267c9-73b8-4ce6-a489-9e2d96d2ccan",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [access control returns err: CheckCanShowDelayedNewsOrFail]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				id: "b32267c9-73b8-4ce6-a489-9e2d96d2ccag",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				newsFilesService: tt.fields.newsFilesService,
				errorFormatter:   tt.fields.errorFormatter,
				access:           tt.fields.access,
			}
			got, err := n.GetAttach(tt.args.id, tt.args.privileges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttach() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttach() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSearchParams(t *testing.T) {
	var tags []string
	var emptyTags []string

	tags = append(tags, "победа")
	tags = append(tags, "красота")

	type args struct {
		dtoQuery string
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantTags  []string
	}{
		{
			name: "Case terms len > 2",
			args: args{
				dtoQuery: "новость #победа #красота",
			},
			wantQuery: "новость",
			wantTags:  tags,
		},
		{
			name: "Case empty query",
			args: args{
				dtoQuery: "",
			},
			wantQuery: "",
			wantTags:  emptyTags,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotTags := getSearchParams(tt.args.dtoQuery)
			if gotQuery != tt.wantQuery {
				t.Errorf("getSearchParams() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("getSearchParams() gotTags = %v, want %v", gotTags, tt.wantTags)
			}
		})
	}
}

func TestNewsService_GetOneBySlug(t *testing.T) {

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	newsRepository := &newsMock.NewsRepositoryMock{}
	fileService := NewNewsFileService(fileRepository, fileStorageRepository)

	type fields struct {
		newsRepository   repository.NewsRepository
		tagRepository    tagrepository.TagRepository
		newsFilesService *NewsFileService
		cache            cache.CacheInterface
		errorFormatter   *helpers.ErrorFormatter
		hash             *helpers.Hash
		access           *service.NewsAccess
	}
	type args struct {
		slug       string
		privileges dto.Privileges
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.News
		wantErr bool
	}{
		{
			name: "Case [no privileges]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
			},
			args: args{
				slug: "123",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [success answer]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
			},
			args: args{
				slug: "6c1729bc-b87b-47d5-b593-d7090a062f37",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_news_detail"},
					},
				},
			},
			want: &entity.News{
				Slug:   "6c1729bc-b87b-47d5-b593-d7090a062f37",
				Active: true,
			},
			wantErr: false,
		},
		{
			name: "Case [getRecord returns err]",
			fields: fields{
				newsRepository:   newsRepository,
				newsFilesService: fileService,
				access:           service.NewNewsAccess(),
				errorFormatter:   &helpers.ErrorFormatter{},
			},
			args: args{
				slug: "6c1729bc-b87b-47d5-b593-d7090a062f40",
				privileges: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_news_detail"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewsService{
				newsRepository:   tt.fields.newsRepository,
				tagRepository:    tt.fields.tagRepository,
				newsFilesService: tt.fields.newsFilesService,
				cache:            tt.fields.cache,
				errorFormatter:   tt.fields.errorFormatter,
				hash:             tt.fields.hash,
				access:           tt.fields.access,
			}
			got, err := n.GetOneBySlug(tt.args.slug, tt.args.privileges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOneBySlug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOneBySlug() got = %v, want %v", got, tt.want)
			}
		})
	}
}
