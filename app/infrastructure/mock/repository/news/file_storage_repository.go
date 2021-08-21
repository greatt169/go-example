package news

import (
	"errors"
	"news-ms/domain/news/entity"
	"time"
)

type FileStorageRepositoryMock struct {
}

// Получить файл по ID
func (f *FileStorageRepositoryMock) GetLink(fileData *entity.File, second time.Duration) (string, error) {
	link := ""
	if fileData.ID == "c32267c9-73b8-4ce6-a489-9e2d96d2ccab" {
		link = "/content-files/news/6c1729bc-b87b-47d5-b593-d7090a062f37/1597771527_Суперновость, файл 1?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=123qwe123qwe123qwe%2F20200907%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200907T125915Z&X-Amz-Expires=60&X-Amz-SignedHeaders=host&response-Content-Disposition=attachment%3B%20filename%3D%22%D0%A1%D1%83%D0%BF%D0%B5%D1%80%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D1%8C%2C%20%D1%84%D0%B0%D0%B9%D0%BB%201.txt%22&response-Content-Type=text%2Fplain%3B%20charset%3Dutf-8&X-Amz-Signature=436b8d7973c13bad5467a64c682d3c7cd8c66e758d176b520f2a3278d49888db"
	}
	if fileData.ID == "c32267c9-73b8-4ce6-a489-9e2d96d2ccam" {
		return "", errors.New("fail_storage")
	}
	return link, nil
}

// Сохранение файлов из сущости новостей
func (f *FileStorageRepositoryMock) SaveNewsFiles(news *entity.News) ([]entity.File, error) {
	var files []entity.File
	// если нет uid
	if len(news.Files) > 0 {
		if news.Files[0].Name == "success" {
			return news.Files, nil
		}
	}
	if news.Id.String() == "c32267c9-73b8-4ce6-a489-9e2d96d2ccab" || news.Id.String() == "c32267c9-73b8-4ce6-a489-9e2d96d2ccad" || news.Id.String() == "6c1729bc-b87b-47d5-b593-d7090a062f37" {
		return news.Files, nil
	} else {
		return files, errors.New("fail1")
	}
}

// Удаляет не переданные в запросе на обновление файлы
func (f *FileStorageRepositoryMock) removeNotPassedNewsFiles(news *entity.News) error {
	return nil
}

// Удаление файлов внутри директории
func (f FileStorageRepositoryMock) RemoveFolder(folderName string) error {
	return nil
}

// Удаление файла
func (f *FileStorageRepositoryMock) RemoveFile(filePath string) error {
	return nil
}
