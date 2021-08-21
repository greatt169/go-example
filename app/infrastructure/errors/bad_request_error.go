package errors

import (
	"fmt"
	"net/http"
)

//Внутренняя ошибка сервиса, определенная на уровне инфраструктуры
type BadRequestError struct {
	Err     error
	Message string
	Trace   string
}

func (e BadRequestError) GetErrorLevel() string {
	return "error"
}
func (e BadRequestError) GetAppErrorCode() string {
	return "BAD_REQUEST"
}

func (e BadRequestError) GetAppErrorStatus() int {
	return http.StatusBadRequest
}

func (e BadRequestError) GetAppErrorMessage() string {
	return "Некорректный запрос, отсутствует один из обязательных параметров."
}

func (e BadRequestError) GetErrorDebugInfo() string {
	return fmt.Sprintf("Geo service error: %v: %v", e.Message, e.Error())
}

func (e BadRequestError) Error() string {
	return e.Err.Error()
}
func (e BadRequestError) GetError() error {
	return e.Err
}
