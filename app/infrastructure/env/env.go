package env

import (
	"fmt"
	"os"
)

// Format - this format text duplicated in code sometimes
const Format = "/%s"

// All variables for project
var (
	Dbdriver             = Getter("DB_DRIVER", "")
	Host                 = Getter("DB_HOST", "")
	Password             = Getter("DB_PASSWORD", "")
	User                 = Getter("DB_USER", "")
	Dbname               = Getter("DB_NAME", "")
	DbPort               = Getter("DB_PORT", "")
	GrpcPort             = fmt.Sprintf(":%s", Getter("GRPC_PORT", "50051"))
	Port                 = fmt.Sprintf(":%s", Getter("SERVER_PORT", "8080"))
	ServiceName          = Getter("API_PATH_SERVICE_NAME", "news-ms")
	PrometheusPort       = fmt.Sprintf("%s", Getter("PROMETHEUS_PORT", "9092"))
	CredentialsFromVault = Getter("CREDENTIALS_FROM_VAULT", "false")
	VaultUrl             = Getter("VAULT_URL", "")
	VaultMountPoint      = Getter("VAULT_MOUNT_POINT", "")
	VaultSecretId        = Getter("VAULT_SECRET_ID", "")
	VaultRoleId          = Getter("VAULT_ROLE_ID", "")
	VaultSecretsBasePath = Getter("VAULT_SECRETS_BASE_PATH", "")
	S3Endpoint           = Getter("S3_ENDPOINT", "")
	S3AccessKey          = Getter("S3_ACCESS_KEY", "")
	S3Secret             = Getter("S3_SECRET", "")
	S3SecureConnection   = Getter("S3_SECURE_MODE", "false")
	S3Region             = Getter("S3_REGION", "")
	S3TraceON            = Getter("S3_LOGMODE", "")
	S3Bucket             = Getter("S3_BUCKET", "")
)

// Getter -
func Getter(key, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	return defaultValue
}
