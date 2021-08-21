package service

import (
	"fmt"
	"geo/infrastructure/interfaces/errors"
	"geo/interfaces/rest/dto"
	"github.com/gin-gonic/gin"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const unexpectedErrorCode = "INTERNAL_SERVER_ERROR"

// Сервис обработки ошибок
type Error struct {
	Logger log.FieldLogger
}

func NewError(logger log.FieldLogger) *Error {
	return &Error{
		Logger: logger,
	}
}

func (e Error) GetResponseError(err error) dto.Error {
	errorFormatted := dto.Error{}

	handledError, ok := err.(errors.HandledError)
	if ok {
		errorFormatted.Error = handledError.GetAppErrorCode()
		errorFormatted.Code = handledError.GetAppErrorStatus()
		errorFormatted.Message = handledError.GetAppErrorMessage()
		errorFormatted.Debug = handledError.GetErrorDebugInfo()
	} else {
		errorFormatted.Error = unexpectedErrorCode
		errorFormatted.Code = http.StatusInternalServerError
		errorFormatted.Message = "Произошла ошибка. Попробуйте выполнить операцию позже"
		errorFormatted.Debug = err.Error()
	}

	return errorFormatted
}

type stackTracer interface {
	StackTrace() pkgErrors.StackTrace
}

func (e Error) LogError(err error) {

	handledError, ok := err.(errors.HandledError)
	if ok {
		switch handledError.GetErrorLevel() {
		case "error":
			tracedErr, ok := handledError.GetError().(stackTracer)
			if !ok {
				panic("Error does not implement stackTracer")
			}

			st := tracedErr.StackTrace()
			e.Logger.WithFields(log.Fields{
				"stackTrace": fmt.Errorf("%+v", st[0:]),
			}).Error(fmt.Sprintf("geo service error: %v", err.Error()))

		case "info":
			tracedErr, ok := handledError.GetError().(stackTracer)
			if !ok {
				panic("Error does not implement stackTracer")
			}

			st := tracedErr.StackTrace()
			e.Logger.WithFields(log.Fields{
				"stackTrace": fmt.Errorf("%+v", st[0:]),
			}).Info(fmt.Sprintf("geo service query error: %v", err.Error()))
		}
	} else {
		e.Logger.Error(fmt.Sprintf("geo service unexpected error: %v", err.Error()))
	}

}

func (e Error) HandleError(err error, c *gin.Context) {
	e.LogError(err)
	errorFormatted := e.GetResponseError(err)
	code := errorFormatted.Code
	c.JSON(code, errorFormatted)
}
