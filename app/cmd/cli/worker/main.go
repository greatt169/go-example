package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"news-ms/infrastructure/registry"
	"os"
)

func main() {
	ctn, err := registry.NewContainer()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}
	// Получение лог-сервиса из контейнера
	logger := ctn.Resolve("logger").(log.FieldLogger)
	if err != nil {
		logger.Error(err.Error())
	}
	// Получение лог-сервиса из контейнера
	app := ctn.Resolve("cli").(*cli.App)
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
