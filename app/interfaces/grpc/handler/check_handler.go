package handler

import (
	"context"
	"news-ms/infrastructure/env"
	"news-ms/interfaces/grpc/proto/v1/news"
)

// Обработчик запросов сервиса
type CheckHandler struct{}

// Конструктор
func NewCheckHandler() *CheckHandler {
	return &CheckHandler{}
}

// Проверка состояния сервиса
func (c CheckHandler) CheckHealth(context.Context, *content_v1.EmptyRequest) (*content_v1.HealthResponse, error) {
	return &content_v1.HealthResponse{
		ServiceName:   env.ServiceName,
		ServiceStatus: "UP",
	}, nil
}
