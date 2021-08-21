package news

import (
	helpers "github.com/AeroAgency/golang-helpers-lib"
	uuid "github.com/satori/go.uuid"
	"news-ms/application/dto/news"
	"news-ms/application/service"
	"news-ms/domain/news/entity"
	"news-ms/domain/news/repository"
	newsMock "news-ms/infrastructure/mock/repository/news"
	"news-ms/infrastructure/persistence/awss3"
	"news-ms/infrastructure/persistence/postgres"
	"reflect"
	"testing"
)

func TestNewNewsFileService(t *testing.T) {
	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	NewNewsFileService(fileRepository, fileStorageRepository)
}

func TestNewsFileService_checkFilesNewsOrFail(t *testing.T) {

	fileRepository := postgres.FileRepository{}
	fileStorageRepository := awss3.FileStorageRepository{}
	validatorRules := service.ValidatorRulesMap{}
	rulesWithSizeMoreThanMax := map[string]string{
		"file_ext":    "doc,xls,ppt,jpg,bmp,pdf,rtf,txt,zip",
		"file_size":   "1",
		"files_limit": "7",
	}

	type fields struct {
		fileRepository        repository.FileRepository
		fileStorageRepository repository.FileStorageInterface
		errorFormatter        *helpers.ErrorFormatter
		fileChecker           *helpers.FileChecker
	}
	type args struct {
		newsFiles []entity.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Check files. Case [valid_files]",
			fields: fields{
				fileRepository:        &fileRepository,
				fileStorageRepository: &fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				newsFiles: []entity.File{
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Check files. Case [count file more than max]",
			fields: fields{
				fileRepository:        &fileRepository,
				fileStorageRepository: &fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           helpers.NewFileChecker(validatorRules.GetNewsFilesMap()),
			},
			args: args{
				newsFiles: []entity.File{
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Check files. Case [invalid file ext]",
			fields: fields{
				fileRepository:        &fileRepository,
				fileStorageRepository: &fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           helpers.NewFileChecker(validatorRules.GetNewsFilesMap()),
			},
			args: args{
				newsFiles: []entity.File{
					{
						Name:   "Суперновость, файл 1",
						Ext:    "invalid",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Check files. Case [invalid file id]",
			fields: fields{
				fileRepository:        &fileRepository,
				fileStorageRepository: &fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           helpers.NewFileChecker(validatorRules.GetNewsFilesMap()),
			},
			args: args{
				newsFiles: []entity.File{
					{
						ID:     "invalid",
						Name:   "Суперновость, файл 1",
						Ext:    "invalid",
						Base64: "SGVsbG8gV29ybGQK",
					},
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Check files. Case [size more than max]",
			fields: fields{
				fileRepository:        &fileRepository,
				fileStorageRepository: &fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           helpers.NewFileChecker(rulesWithSizeMoreThanMax),
			},
			args: args{
				newsFiles: []entity.File{
					{
						Name:   "Суперновость, файл 1",
						Ext:    "txt",
						Base64: "SGVsbG8gV29ybGQK",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewsFileService{
				fileRepository:        tt.fields.fileRepository,
				fileStorageRepository: tt.fields.fileStorageRepository,
				errorFormatter:        tt.fields.errorFormatter,
				fileChecker:           tt.fields.fileChecker,
			}
			if err := f.checkFilesNewsOrFail(tt.args.newsFiles); (err != nil) != tt.wantErr {
				t.Errorf("checkFilesNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsFileService_GetFileInfo(t *testing.T) {
	type fields struct {
		fileRepository        repository.FileRepository
		fileStorageRepository repository.FileStorageInterface
		errorFormatter        *helpers.ErrorFormatter
		fileChecker           *helpers.FileChecker
	}
	type args struct {
		id string
	}

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *news.FileInfoDto
		wantErr bool
	}{
		{
			name: "Case [Repo answers ok, Storage answers ok]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
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
			name: "Case [Repo answers err, Storage answers -]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				id: "c32267c9-73b8-4ce6-a489-9e2d96d2ccac",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case [Repo answers ok, Storage answers err]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				id: "c32267c9-73b8-4ce6-a489-9e2d96d2ccam",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewsFileService{
				fileRepository:        tt.fields.fileRepository,
				fileStorageRepository: tt.fields.fileStorageRepository,
				errorFormatter:        tt.fields.errorFormatter,
				fileChecker:           tt.fields.fileChecker,
			}
			got, err := f.GetFileInfo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsFileService_DeleteFiles(t *testing.T) {
	type fields struct {
		fileRepository        repository.FileRepository
		fileStorageRepository repository.FileStorageInterface
		errorFormatter        *helpers.ErrorFormatter
		fileChecker           *helpers.FileChecker
	}
	type args struct {
		id uuid.UUID
	}

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case [Repo answers ok]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccab"),
			},
			wantErr: false,
		},
		{
			name: "Case [Repo answers err]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccam"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewsFileService{
				fileRepository:        tt.fields.fileRepository,
				fileStorageRepository: tt.fields.fileStorageRepository,
				errorFormatter:        tt.fields.errorFormatter,
				fileChecker:           tt.fields.fileChecker,
			}
			if err := f.DeleteFiles(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsFileService_SaveFiles(t *testing.T) {
	type fields struct {
		fileRepository        repository.FileRepository
		fileStorageRepository repository.FileStorageInterface
		errorFormatter        *helpers.ErrorFormatter
		fileChecker           *helpers.FileChecker
	}
	type args struct {
		newsEntity entity.News
	}

	fileRepository := &newsMock.FileRepositoryMock{}
	fileStorageRepository := &newsMock.FileStorageRepositoryMock{}
	validatorRules := service.ValidatorRulesMap{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case [Invalid files]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           helpers.NewFileChecker(validatorRules.GetNewsFilesMap()),
			},
			args: args{
				newsEntity: entity.News{
					Id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccab"),
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
			name: "Case [Repo answers ok, Storage answers ok]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				newsEntity: entity.News{
					Id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccab"),
					Files: []entity.File{
						{
							Name:   "Суперновость, файл 1",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Case [Storage answers err,  Repo answers -]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				newsEntity: entity.News{
					Id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccac"),
					Files: []entity.File{
						{
							Name:   "Суперновость, файл 1",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Case [Storage answers ok,  Repo answers err]",
			fields: fields{
				fileRepository:        fileRepository,
				fileStorageRepository: fileStorageRepository,
				errorFormatter:        &helpers.ErrorFormatter{},
				fileChecker:           &helpers.FileChecker{},
			},
			args: args{
				newsEntity: entity.News{
					Id: helpers.GetUuidByString("c32267c9-73b8-4ce6-a489-9e2d96d2ccad"),
					Files: []entity.File{
						{
							Name:   "Суперновость, файл 1",
							Ext:    "txt",
							Base64: "SGVsbG8gV29ybGQK",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewsFileService{
				fileRepository:        tt.fields.fileRepository,
				fileStorageRepository: tt.fields.fileStorageRepository,
				errorFormatter:        tt.fields.errorFormatter,
				fileChecker:           tt.fields.fileChecker,
			}
			if err := f.SaveFiles(tt.args.newsEntity); (err != nil) != tt.wantErr {
				t.Errorf("SaveRequestFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
