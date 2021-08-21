package registry

import (
	"errors"
	"fmt"
	cache "github.com/AeroAgency/golang-bigcache-lib"
	filestorage "github.com/AeroAgency/golang-filestorage-lib"
	"github.com/allegro/bigcache"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/credentials"
	"github.com/sarulabs/di"
	log "github.com/sirupsen/logrus"
	"news-ms/application/usecases/news"
	newsrepository "news-ms/domain/news/repository"
	tagrepository "news-ms/domain/tag/repository"
	"news-ms/infrastructure/env"
	"news-ms/infrastructure/grpc"
	awss3Repository "news-ms/infrastructure/persistence/awss3"
	"news-ms/infrastructure/persistence/postgres"
	"news-ms/infrastructure/vars"
	"news-ms/infrastructure/vault"
	"os"
	"strconv"
	"time"
)

// Контейнер внедрения зависимостей
type Container struct {
	ctn di.Container
}

// Конструктор для контейнера
func NewContainer() (*Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}
	if err := builder.Add([]di.Def{
		{ // grpc-server
			Name:  "grpc_server",
			Build: buildGrpcServerContainer(),
		},
		{ // Логгер
			Name:  "logger",
			Build: buildLoggerContainer(),
		},
		{ // Модель БД
			Name:  "db",
			Build: buildDataBaseContainer,
			Close: func(obj interface{}) error {
				return obj.(*postgres.Database).Close()
			},
		},
		{ // Репозиторий новостей
			Name: "news_repository",
			Build: func(ctn di.Container) (interface{}, error) {
				dbContainer := ctn.Get("db").(*postgres.Database)
				return postgres.NewNewsRepository(dbContainer.DB), nil
			},
		},
		{ // Репозиторий тегов
			Name: "tag_repository",
			Build: func(ctn di.Container) (interface{}, error) {
				dbContainer := ctn.Get("db").(*postgres.Database)
				return postgres.NewTagRepository(dbContainer.DB), nil
			},
		},
		{ // Репозиторий файлов (бд)
			Name: "file_repository",
			Build: func(ctn di.Container) (interface{}, error) {
				dbContainer := ctn.Get("db").(*postgres.Database)
				return postgres.NewFileRepository(dbContainer.DB), nil
			},
		},
		{ // Репозиторий файлов (файловое хранилище)
			Name:  "file_storage_repository",
			Build: buildFileStorageContainer,
		},
		{ // Модель кэша
			Name:  "cache",
			Build: buildCacheContainer,
		},
		{ // Сервис файлов новостей
			Name:  "news_files_service",
			Build: buildNewsFilesServiceContainer(),
		},
		{ // Сервис новостей
			Name:  "news_service",
			Build: buildNewsServiceContainer(),
		},
		{ // Модель работы с vault
			Name:  "vault_client",
			Build: buildVaultClientContainer,
		},
		{ // Модель работы с переменными
			Name:  "vars",
			Build: buildVarsContainer,
		},
	}...); err != nil {
		return nil, err
	}

	return &Container{
		ctn: builder.Build(),
	}, nil
}

// собирает контейнер для работы с клиентом vault
func buildVaultClientContainer(ctn di.Container) (interface{}, error) {
	// Получение лог-сервиса из контейнера
	logger := ctn.Get("logger").(log.FieldLogger)
	vaultClient := vault.NewVaultClient(logger)
	return vaultClient, nil
}

// собирает контейнер для работы с клиентом vault
func buildVarsContainer(ctn di.Container) (interface{}, error) {
	// Получение vault-сервиса из контейнера
	vaultClient := ctn.Get("vault_client").(*vault.VaultClient)
	vars := vars.NewVarsModel(vaultClient)
	return vars, nil
}

//Собирает контейнер для логов
func buildLoggerContainer() func(ctn di.Container) (interface{}, error) {
	return func(ctn di.Container) (interface{}, error) {
		logger := log.Logger{}
		logger.SetFormatter(&log.JSONFormatter{})
		logger.SetOutput(os.Stdout)
		logger.Info("run grpc server")
		logger.SetLevel(log.DebugLevel)
		return &logger, nil
	}
}

//Собирает контейнер для grpc-сервера
func buildGrpcServerContainer() func(ctn di.Container) (interface{}, error) {
	return func(ctn di.Container) (interface{}, error) {
		// Получение лог-сервиса из контейнера
		logger := ctn.Get("logger").(log.FieldLogger)
		logger.Info("Try to build grpc server")
		// Получение сервиса для работы с новостями из контейнера
		newsService := ctn.Get("news_service").(*news.NewsService)
		// Получение объекта для работы с кэшем из контейнера
		cache := ctn.Get("cache").(cache.CacheInterface)
		// Получение репозитория тегов
		tagRepository := ctn.Get("tag_repository").(tagrepository.TagRepository)
		server := grpc.NewGrpcServer(
			logger,
			newsService,
			cache,
			tagRepository,
		)
		logger.Info("grpc server has been successfully built")
		return server, nil
	}
}

//Собирает контейнер для сервиса новостей
func buildNewsServiceContainer() func(ctn di.Container) (interface{}, error) {
	return func(ctn di.Container) (interface{}, error) {
		newsRepository := ctn.Get("news_repository").(newsrepository.NewsRepository)
		tagRepository := ctn.Get("tag_repository").(tagrepository.TagRepository)
		newsFilesService := ctn.Get("news_files_service").(*news.NewsFileService)
		cache := ctn.Get("cache").(cache.CacheInterface)
		return news.NewNewsService(
			newsRepository,
			tagRepository,
			newsFilesService,
			cache,
		), nil
	}
}

//Собирает контейнер для сервиса файлов новостей
func buildNewsFilesServiceContainer() func(ctn di.Container) (interface{}, error) {
	return func(ctn di.Container) (interface{}, error) {
		fileRepository := ctn.Get("file_repository").(newsrepository.FileRepository)
		fileStorageRepository := ctn.Get("file_storage_repository").(newsrepository.FileStorageInterface)
		return news.NewNewsFileService(
			fileRepository,
			fileStorageRepository,
		), nil
	}
}

// Собирает контейнер с репозиториями БД
func buildDataBaseContainer(ctn di.Container) (interface{}, error) {
	// Получение лог-сервиса из контейнера
	logger := ctn.Get("logger").(log.FieldLogger)
	logger.Info("Try to connect to database")
	var dsn string
	// Получение сервиса для работы с переменными
	vars := ctn.Get("vars").(*vars.VarsModel)
	dsn = "host=" + vars.GetDbHost() + " port=" + vars.GetDbPort() + " user=" + vars.GetDbUser() + " dbname=" + vars.GetDbName() + " sslmode=disable password=" + vars.GetDbPassword()
	repos, err := postgres.NewDatabase(env.Dbdriver, dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("DB error connection. Details: %s", err))
		conn := "host=" + vars.GetDbHost() + " port=" + vars.GetDbPort() + " user=" + vars.GetDbUser() + " dbname=" + vars.GetDbName() + " sslmode=disable"
		logger.Error(fmt.Sprintf("DB error connection. Connection data: %s", conn))
		return &postgres.Database{DB: nil}, nil
	}
	logger.Info("DB connection has  been successfully created")
	return repos, err
}

// Собирает контейнер для кэша
func buildCacheContainer(ctn di.Container) (interface{}, error) {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	bigcacheInstance, _ := bigcache.NewBigCache(config)
	cacheInstance := cache.NewBigCache(bigcacheInstance)
	return cacheInstance, nil
}

// Собирает контейнер для работы с хранилищем (any s3)
func buildFileStorageContainer(ctn di.Container) (i interface{}, err error) {
	var storageSecureConnection bool
	storageSecureConnection, _ = strconv.ParseBool(env.S3SecureConnection)
	// Получение лог-сервиса из контейнера
	logger := ctn.Get("logger").(log.FieldLogger)
	logger.Info("Try to connect to file storage")
	// Получение сервиса для работы с переменными
	vars := ctn.Get("vars").(*vars.VarsModel)
	c1 := make(chan newsrepository.FileStorageInterface, 1)
	// Run your long running function in it's own goroutine and pass back it's
	// response into our channel.
	go func() {
		s3Client, err := minio.NewWithOptions(env.S3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(vars.GetS3AccessKey(), vars.GetS3Secret(), ""),
			Secure: storageSecureConnection,
			Region: env.S3Region,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("s3 error while getting client. Details: %s", err))
			c1 <- nil
			return
		}
		logmode, err := strconv.Atoi(env.S3TraceON)
		if err != nil {
			logmode = 0
		}
		if logmode == 1 {
			s3Client.TraceOn(os.Stdout)
		}
		storage, err := filestorage.NewMinioFileStorage(s3Client, env.S3Bucket)
		if err != nil {
			logger.Error(fmt.Sprintf("s3 error connection. Details: %s", err))
			c1 <- nil
			return
		}
		storageRepository := awss3Repository.NewFileStorageRepository(storage)
		logger.Info("File storage connection has  been successfully created")
		c1 <- storageRepository
	}()
	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case storageRepository := <-c1:
		if storageRepository == nil {
			storageRepository = &awss3Repository.FileStorageRepository{}
		}
		return storageRepository, nil
	case <-time.After(5 * time.Second):
		err = errors.New("5 seconds timeout error")
		logger.Error(fmt.Sprintf("s3 error connection. Details: %s", err))
		return &awss3Repository.FileStorageRepository{}, nil
	}
}

// Получение зависимости из контейнера
func (c *Container) Resolve(name string) interface{} {
	return c.ctn.Get(name)
}

// Очистка контейнера
func (c *Container) Clean() error {
	return c.ctn.Clean()
}
