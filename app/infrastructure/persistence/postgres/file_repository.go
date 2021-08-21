package postgres

import (
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"news-ms/domain/news/entity"
)

// Репозиторий для работы с файлами (таблица в БД)
type FileRepository struct {
	db *gorm.DB
}

// Конструктор
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

// Создать файл
func (f *FileRepository) Create(file *entity.File) error {
	f.db.Create(file)
	return f.db.Error
}

// Удалить файл
func (f *FileRepository) Delete(id string) error {
	file := entity.File{ID: id}
	db := f.db.Unscoped().Delete(&file)
	return db.Error
}

// Получить файл по ID
func (f *FileRepository) GetByID(id string) (*entity.File, error) {
	var file entity.File
	db := f.db.Where("id=?", id).
		Find(&file)
	if db.Error != nil {
		return nil, db.Error
	}
	return &file, db.Error
}

// Получить файлы по entityId
func (f *FileRepository) getByEntityID(uuid uuid.UUID) ([]entity.File, error) {
	var files []entity.File
	db := f.db.Where("entity_id=?", uuid).
		Find(&files)
	if db.Error != nil {
		return nil, db.Error
	}
	return files, db.Error
}

// Сохранить переданные файлы
func (f *FileRepository) savePassedNewsFiles(files []entity.File) error {
	for _, v := range files {
		file, _ := f.GetByID(v.ID)
		if file != nil {
			continue
		}
		err := f.Create(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Удалить не переданные файлы для новости
func (f *FileRepository) removeNotPassedNewsFiles(entityUuid uuid.UUID, passedFiles []entity.File) error {
	var passedFileIds []string
	existedFiles, err := f.getByEntityID(entityUuid)
	if err != nil {
		return err
	}
	for _, passedFile := range passedFiles {
		passedFileIds = append(passedFileIds, passedFile.ID)
	}
	for _, existedFile := range existedFiles {
		res := helpers.StringInSlice(existedFile.ID, passedFileIds)
		if res != true {
			err = f.Delete(existedFile.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Сохраняет переданные файлы, не переданные - удаляет
func (f *FileRepository) SaveNewsFiles(entityUuid uuid.UUID, files []entity.File) error {
	err := f.removeNotPassedNewsFiles(entityUuid, files)
	if err != nil {
		return err
	}
	err = f.savePassedNewsFiles(files)
	if err != nil {
		return err
	}
	return nil
}

// Удалить все файлы
func (f *FileRepository) RemoveAll() error {
	f.db.Delete(entity.File{})
	return nil
}
