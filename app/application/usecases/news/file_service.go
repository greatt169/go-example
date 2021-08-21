package news

import (
	"errors"
	"fmt"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"news-ms/application/dto/news"
	"news-ms/application/service"
	"news-ms/domain/news/entity"
	"news-ms/domain/news/repository"
	"time"
)

// Сервис для работы с файлами новостей
type NewsFileService struct {
	fileRepository        repository.FileRepository
	fileStorageRepository repository.FileStorageInterface
	errorFormatter        *helpers.ErrorFormatter
	fileChecker           *helpers.FileChecker
}

// Конструктор
func NewNewsFileService(
	fileRepository repository.FileRepository,
	fileStorageRepository repository.FileStorageInterface,
) *NewsFileService {
	validatorRules := &service.ValidatorRulesMap{}
	return &NewsFileService{
		fileRepository:        fileRepository,
		fileStorageRepository: fileStorageRepository,
		errorFormatter:        &helpers.ErrorFormatter{},
		fileChecker:           helpers.NewFileChecker(validatorRules.GetNewsFilesMap()),
	}
}

// Удаляет файлы новости из хранилища
func (f NewsFileService) DeleteFiles(id uuid.UUID) error {
	// Удаляем файл из хранилища
	_, err := f.fileStorageRepository.SaveNewsFiles(&entity.News{Id: id})
	if err != nil {
		return err
	}
	return nil
}

// Сохраняет файлы новости
func (f NewsFileService) SaveFiles(newsEntity entity.News) error {
	// Проверка файлов на корректность
	err := f.checkFilesNewsOrFail(newsEntity.Files)
	if err != nil {
		return err
	}
	// Сохраняем файлы новости в хранилище
	filesInfo, err := f.fileStorageRepository.SaveNewsFiles(&newsEntity)
	if err != nil {
		return err
	}
	// Сохраняем информацию о файлах в БД
	err = f.fileRepository.SaveNewsFiles(newsEntity.Id, filesInfo)
	if err != nil {
		return f.errorFormatter.ReturnError(codes.Internal, err, "Ошибка создания записей о файлах в БД")
	}
	return nil
}

// Возвращает информацию о файле новости
func (f NewsFileService) GetFileInfo(id string) (*news.FileInfoDto, error) {
	//Получаем информацию о файле из БД
	fileData, err := f.fileRepository.GetByID(id)
	if err != nil {
		return nil, f.errorFormatter.ReturnError(codes.NotFound, err, "")
	}
	//Получаем ссылку на файл
	link, err := f.fileStorageRepository.GetLink(fileData, 60*time.Second)
	if err != nil {
		err = f.errorFormatter.Wrap(err, err.Error())
		return nil, f.errorFormatter.ReturnError(codes.Internal, err, "")
	}
	return &news.FileInfoDto{
		Fields: fileData,
		Link:   link,
	}, nil
}

// Проверяет вложения новостей на корректность
func (f NewsFileService) checkFilesNewsOrFail(newsFiles []entity.File) error {
	isValidCount := f.fileChecker.IsValidCount(len(newsFiles))
	if isValidCount != true {
		maxFileCount := f.fileChecker.GetMaxFilesCount()
		errorMessage := fmt.Sprintf("Можно прикрепить только %d файлов.", maxFileCount)
		return f.errorFormatter.ReturnError(codes.FailedPrecondition, errors.New(""), errorMessage)
	}
	for _, newsFile := range newsFiles {
		if newsFile.ID != "" {
			if f.fileChecker.IsValidFileId(newsFile.ID) != true {
				errorText := "File (" + newsFile.Name + "." + newsFile.Ext + ") couldn't be uploaded, because id (" + newsFile.ID + ") doesn't match with uuid pattern."
				errorMessage := "Недопустимый формат Id файла"
				return f.errorFormatter.ReturnError(codes.FailedPrecondition, errors.New(errorText), errorMessage)
			}
		}
		isValidExt := f.fileChecker.IsValidExt(newsFile.Ext)
		if isValidExt != true {
			errorText := "File (" + newsFile.Name + "." + newsFile.Ext + ") couldn't be uploaded, because extension (" + newsFile.Ext + ") doesn't support."
			errorMessage := "Недопустимый формат файла"
			return f.errorFormatter.ReturnError(codes.FailedPrecondition, errors.New(errorText), errorMessage)
		}
		isValidSize := f.fileChecker.IsValidSize(newsFile.Base64)
		if isValidSize != true {
			errorText := "File (" + newsFile.Name + "." + newsFile.Ext + ") couldn't be uploaded, because size is bigger then max."
			errorMessage := "Превышен допустимый размер файла"
			return f.errorFormatter.ReturnError(codes.FailedPrecondition, errors.New(errorText), errorMessage)
		}
	}
	return nil
}
