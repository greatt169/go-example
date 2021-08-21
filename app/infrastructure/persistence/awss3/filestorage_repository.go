package awss3

import (
	"bytes"
	"encoding/base64"
	filestorage "github.com/AeroAgency/golang-filestorage-lib"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"news-ms/domain/news/entity"
	"strings"
	"time"
)

type FileStorageRepository struct {
	fileStorage filestorage.FileStorageInterface
}

// Конструктор
func NewFileStorageRepository(storageInterface filestorage.FileStorageInterface) *FileStorageRepository {
	return &FileStorageRepository{
		fileStorage: storageInterface,
	}
}

// Возвращает путь к файлу
func getNewsFilePath(entityId uuid.UUID, fileId string, file entity.File) string {
	filePath := "news/" + entityId.String() + "/" + fileId + "/" + file.Name + "." + file.Ext
	return filePath
}

// Сохранение файлов из сущости новостей
func (f *FileStorageRepository) SaveNewsFiles(news *entity.News) ([]entity.File, error) {
	// Удаляем не переданные файлы
	err := f.removeNotPassedNewsFiles(news)
	if err != nil {
		return nil, err
	}
	var files []entity.File
	var fileId string
	errorFormatter := helpers.ErrorFormatter{}
	for _, file := range news.Files {
		if file.Name != "" {
			if file.ID != "" {
				fileId = file.ID
			} else {
				fileId = uuid.NewV4().String()
			}
			// Обработка файлов для загрузки
			timeCreateUnix := time.Now().Unix()
			filePath := getNewsFilePath(news.Id, fileId, file)
			dec, err := base64.StdEncoding.DecodeString(file.Base64[strings.IndexByte(file.Base64, ',')+1:])
			if err != nil {
				err = errorFormatter.Wrap(err, "Caught error while decoding file from base64")
				err = errorFormatter.ReturnError(codes.Internal, err, "")
				return nil, err
			}
			// Загрузка в хранилище
			err = f.fileStorage.UploadFile(f.fileStorage.GetBucketName(), filePath, bytes.NewReader(dec), int64(len(dec)))
			if err != nil {
				err = errorFormatter.Wrap(err, "Caught error while uploading file to storage")
				err = errorFormatter.ReturnError(codes.Internal, err, "")
				return nil, err
			}
			// Формирование ответа в виде файлов для отображения
			file := entity.File{
				ID:         fileId,
				DateCreate: timeCreateUnix,
				Name:       file.Name,
				Ext:        file.Ext,
				Path:       filePath,
				UserId:     news.UserId,
				EntityId:   news.Id,
				Type:       "file",
			}
			files = append(files, file)
		}
	}
	return files, nil
}

// Удаляет не переданные в запросе на обновление файлы
func (f *FileStorageRepository) removeNotPassedNewsFiles(news *entity.News) error {
	var passedFilePaths []string
	for _, file := range news.Files {
		if file.ID != "" {
			passedFilePaths = append(passedFilePaths, getNewsFilePath(news.Id, file.ID, file))
		}
	}
	newsFolderPath := "news" + "/" + news.Id.String()
	existedFilePaths, err := f.fileStorage.GetFilesIntoFolder(f.fileStorage.GetBucketName(), newsFolderPath)
	if err != nil {
		return err
	}
	for _, existedFilePath := range existedFilePaths {
		res := helpers.StringInSlice(existedFilePath, passedFilePaths)
		if res != true {
			err = f.RemoveFile(existedFilePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Получение ссылки на файл
func (f *FileStorageRepository) GetLink(fileData *entity.File, second time.Duration) (string, error) {
	fileName := fileData.Name + "." + fileData.Ext
	link, err := f.fileStorage.GetFileLink(f.fileStorage.GetBucketName(), fileName, fileData.Path, second)
	if err != nil {
		return "", err
	}
	return link, nil
}

// Удаление файлов внутри директории
func (f FileStorageRepository) RemoveFolder(folderName string) error {
	err := f.fileStorage.RemoveFolder(f.fileStorage.GetBucketName(), folderName)
	if err != nil {
		return err
	}
	return nil
}

// Удаление файла
func (f *FileStorageRepository) RemoveFile(filePath string) error {
	err := f.fileStorage.RemoveFile(f.fileStorage.GetBucketName(), filePath)
	if err != nil {
		return err
	}
	return nil
}
