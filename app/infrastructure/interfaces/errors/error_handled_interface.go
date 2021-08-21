package errors

type HandledError interface {
	GetAppErrorCode() string
	GetAppErrorStatus() int
	GetAppErrorMessage() string
	GetErrorDebugInfo() string
	GetErrorLevel() string
	GetError() error
}
