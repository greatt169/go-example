package vars

import (
	"news-ms/infrastructure/env"
	"news-ms/infrastructure/vault"
	"strconv"
)

// Модель объекта работы с переменными в 2 режимах:
// - vault
// - env
type VarsModel struct {
	secrets *vault.VaultCredentials
}

func NewVarsModel(vaultClient *vault.VaultClient) *VarsModel {
	secrets, err := vaultClient.GetVaultSecrets()
	if err != nil {
		vaultClient.LogErr(err)
	}
	return &VarsModel{
		secrets: secrets,
	}
}

// Получение хоста БД
func (v VarsModel) GetDbHost() string {
	return env.Host
}

// Получение имени БД
func (v VarsModel) GetDbName() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.DB_NAME
	} else {
		return env.Dbname
	}
}

// Получение пароля БД
func (v VarsModel) GetDbPassword() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.DB_PASSWORD
	} else {
		return env.Password
	}
}

// Получение порта БД
func (v VarsModel) GetDbPort() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.DB_PORT
	} else {
		return env.DbPort
	}
}

// Получение пользователя БД
func (v VarsModel) GetDbUser() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.DB_USER
	} else {
		return env.User
	}
}

// Получение S3AccessKey
func (v VarsModel) GetS3AccessKey() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.S3_ACCESS_KEY
	} else {
		return env.S3AccessKey
	}
}

// Получение GetS3Secret
func (v VarsModel) GetS3Secret() string {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault {
		return v.secrets.S3_SECRET
	} else {
		return env.S3Secret
	}
}
