package news

import (
	"errors"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	uuid "github.com/satori/go.uuid"
	"news-ms/domain/news/entity"
)

type FileRepositoryMock struct {
}

// Получить файл по ID
func (f *FileRepositoryMock) GetByID(id string) (*entity.File, error) {
	var file *entity.File
	if id == "c32267c9-73b8-4ce6-a489-9e2d96d2ccab" || id == "c32267c9-73b8-4ce6-a489-9e2d96d2ccam" {
		file = &entity.File{
			ID:         id,
			Name:       "Суперновость, файл 1",
			Ext:        "txt",
			Path:       "news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1",
			UserId:     "dc72578e-da65-4de2-bbc2-8d3fc82c7ebb",
			DateCreate: 1597771527,
			EntityId:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f37"),
			Type:       "file",
			Base64:     "SGVsbG8gV29ybGQK",
		}
		return file, nil
	}
	if id == "p32267c9-73b8-4ce6-a489-9e2d96d2ccam" {
		file = &entity.File{
			ID:         id,
			Name:       "Суперновость, файл 1",
			Ext:        "txt",
			Path:       "news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1",
			UserId:     "dc72578e-da65-4de2-bbc2-8d3fc82c7ebb",
			DateCreate: 1597771527,
			EntityId:   helpers.GetUuidByString("p32267c9-73b8-4ce6-a489-9e2d96d2ccam"),
			Type:       "file",
			Base64:     "SGVsbG8gV29ybGQK",
		}
		return file, nil
	}

	if id == "k32267c9-73b8-4ce6-a489-9e2d96d2ccan" {
		file = &entity.File{
			ID:         id,
			Name:       "Суперновость, файл 1",
			Ext:        "txt",
			Path:       "news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1",
			UserId:     "dc72578e-da65-4de2-bbc2-8d3fc82c7ebb",
			DateCreate: 1597771527,
			EntityId:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f41"),
			Type:       "file",
			Base64:     "SGVsbG8gV29ybGQK",
		}
		return file, nil
	}

	if id == "b32267c9-73b8-4ce6-a489-9e2d96d2ccag" {
		file = &entity.File{
			ID:         id,
			Name:       "Суперновость, файл 1",
			Ext:        "txt",
			Path:       "news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1",
			UserId:     "dc72578e-da65-4de2-bbc2-8d3fc82c7ebb",
			DateCreate: 1597771527,
			EntityId:   helpers.GetUuidByString("6c1729bc-b87b-47d5-b593-d7090a062f45"),
			Type:       "file",
			Base64:     "SGVsbG8gV29ybGQK",
		}
		return file, nil
	}
	return nil, errors.New("fail2")
}

// Создать файл
func (f *FileRepositoryMock) Create(file *entity.File) error {
	return nil
}

// Удалить файл
func (f *FileRepositoryMock) Delete(id string) error {
	return nil
}

// Получить файлы по entityId
func (f *FileRepositoryMock) getByEntityID(uuid uuid.UUID) ([]entity.File, error) {
	var files []entity.File
	return files, nil
}

// Сохранить переданные файлы
func (f *FileRepositoryMock) savePassedNewsFiles(files []entity.File) error {
	return nil
}

// Удалить не переданные файлы для новости
func (f *FileRepositoryMock) removeNotPassedNewsFiles(entityUuid uuid.UUID, passedFiles []entity.File) error {
	return nil
}

// Сохраняет переданные файлы, не переданные - удаляет
func (f *FileRepositoryMock) SaveNewsFiles(entityUuid uuid.UUID, files []entity.File) error {
	// если нет uid
	if len(files) > 0 {
		if files[0].Name == "success" {
			return nil
		}
	}

	if entityUuid.String() == "c32267c9-73b8-4ce6-a489-9e2d96d2ccab" {
		return nil
	}
	return errors.New("fail3")
}

// Удалить все файлы
func (f *FileRepositoryMock) RemoveAll() error {
	return nil
}
