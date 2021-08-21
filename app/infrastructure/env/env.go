package env

import (
	"fmt"
	"os"
)

// Format - this format text duplicated in code sometimes
const Format = "/%s"

// All variables for project
var (
	Port            = fmt.Sprintf(":%s", Getter("SERVER_PORT", "8080"))
	ServiceName     = Getter("API_PATH_SERVICE_NAME", "geo")
	Root            = Getter("ROOT", "api")
	Version         = Getter("VERSION", "v1")
	TraceHeader     = Getter("TRACE_HEADER", "uber-trace-id")
	DadataApiKey    = Getter("DADATA_API_KEY", "")
	DadataSecretKey = Getter("DADATA_SECRET_KEY", "")
	Dbdriver        = Getter("DB_DRIVER", "")
	DbHost          = Getter("DB_HOST", "")
	DbPassword      = Getter("DB_PASSWORD", "")
	DbUser          = Getter("DB_USER", "")
	DbName          = Getter("DB_NAME", "")
	DbPort          = Getter("DB_PORT", "")
	DbSslMode       = Getter("DB_SSL_MODE", "disable")
	DbSslCertPath   = Getter("DB_SSL_CERT_PATH", "")
	RedisHost       = Getter("REDIS_HOST", "10.7.0.5")
	RedisPort       = Getter("REDIS_PORT", "6379")
	RedisExpTime    = Getter("REDIS_CACHE_EXP_TIME", "15")
)

// Getter -
func Getter(key, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	return defaultValue
}
