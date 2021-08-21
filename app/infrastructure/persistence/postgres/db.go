package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	newsentity "news-ms/domain/news/entity"
	tagentity "news-ms/domain/tag/entity"
	"news-ms/infrastructure/env"
	"os"
	"strconv"
)

// Объект для работы с БД
type Database struct {
	DB *gorm.DB
}

// Создание БД
func createDb(dsn string) (*gorm.DB, error) {
	// Инициализируем соединение
	conn := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable", env.User, env.Password, env.Host)
	// Открываем sql-клиент
	dbClient, err := sql.Open("postgres", conn)
	var db *gorm.DB
	if dbClient != nil {
		// Создаем db
		_, err = dbClient.Exec("create database " + env.Dbname + ";")
		if err != nil {
			return nil, err
		}
		// Открываем соединение
		db, err = gorm.Open(env.Dbdriver, dsn)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("error connection to db host")
	}
	return db, nil
}

// Конструктор
func NewDatabase(Dbdriver, dsn string) (*Database, error) {
	db, err := gorm.Open(Dbdriver, dsn)
	if err != nil {
		// Скорее всего базы нет
		// Пытаемся создать
		db, err = createDb(dsn)
	}
	logmode, err := strconv.Atoi(os.Getenv("DB_LOGMODE"))
	if err != nil {
		logmode = 0
	}
	if logmode == 1 {
		db.LogMode(true)
	}
	return &Database{
		DB: db,
	}, nil
}

// Закрытие соединения
func (s *Database) Close() error {
	return s.DB.Close()
}

// Автомиграция структуры данных
func (s *Database) Automigrate() error {
	newsTags := &newsentity.NewsTags{}
	newsFiles := &newsentity.File{}
	db := s.DB.AutoMigrate(
		newsTags,
		&newsentity.News{},
		&newsentity.File{},
		&tagentity.Tag{},
	)

	// Устанавливаем внешние ключи
	db.Model(newsTags).AddForeignKey("news_id", "news(id)", "CASCADE", "CASCADE")
	db.Model(newsFiles).AddForeignKey("entity_id", "news(id)", "CASCADE", "RESTRICT")

	// Индекс для поиска
	db.DB().Exec("CREATE INDEX IF NOT EXISTS idx_gin_news ON news USING gin ((setweight(to_tsvector('russian', 'name'), 'A') || setweight(to_tsvector('russian', 'text_search'), 'B')))")
	return db.Error
}
