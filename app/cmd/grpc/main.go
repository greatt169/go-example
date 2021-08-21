package main

import (
	log "github.com/sirupsen/logrus"
	"news-ms/infrastructure/grpc"
	"news-ms/infrastructure/persistence/postgres"
	"news-ms/infrastructure/registry"
	"sync"
)

func main() {
	// Инициализация контейнера зависимостей
	ctn, err := registry.NewContainer()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	} else {
		log.Info("Registry container has been successfully build")
	}
	// Получение лог-сервиса из контейнера
	logger := ctn.Resolve("logger").(log.FieldLogger)

	// Получение объекта для работы с БД из контейнера
	log.Infof("Getting database from container")
	db := ctn.Resolve("db").(*postgres.Database)
	if db.DB != nil {
		log.Infof("Starting auto migrates")
		// Запуск автомиграций
		err = db.Automigrate()
		if err != nil {
			logger.Error(err.Error())
		} else {
			log.Infof("auto migrates have been successfully finished")
		}
	}
	// Запуск grpc сервера
	log.Infof("Starting grpc server")
	server := ctn.Resolve("grpc_server").(*grpc.ContentServer)
	wg := sync.WaitGroup{}
	wg.Add(2)
	log.Infof("Running grpc server")
	go server.RunGrpcServer(&wg)
	log.Infof("Running proxy server")
	go server.RunProxyServer(&wg)
	wg.Wait()
	// Очистка контейнера
	ctn.Clean()
}
