package main

import (
	"geo/infrastructure/env"
	"geo/infrastructure/registry"
	"geo/infrastructure/rest"
	gintraceser "github.com/AeroAgency/go-gin-tracer"
	"github.com/opentracing/opentracing-go"
)

func main() {
	tracer, closer, err := gintraceser.SetJaegerTracer(env.TraceHeader)
	if err == nil {
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
	}
	// Инициализация контейнера зависимостей
	ctn, _ := registry.NewContainer()
	server := ctn.Resolve("router").(*rest.Router)
	r := server.Router()
	_ = r.Run(env.Port)
	ctn.Clean()
}
