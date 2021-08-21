package vault

import (
	"encoding/json"
	"fmt"
	"github.com/AeroAgency/golang-helpers-lib/vault"
	log "github.com/sirupsen/logrus"
	"news-ms/infrastructure/env"
	"strconv"
)

// Структура хранения секретов в Vault
type VaultCredentials struct {
	DB_HOST       string
	DB_NAME       string
	DB_PASSWORD   string
	DB_PORT       string
	DB_USER       string
	S3_ACCESS_KEY string
	S3_SECRET     string
}

type VaultClient struct {
	Logger    log.FieldLogger
	LibClient vault.Client
}

func NewVaultClient(logger log.FieldLogger) *VaultClient {
	credentialsFromVault, _ := strconv.ParseBool(env.CredentialsFromVault)
	if credentialsFromVault != true {
		return &VaultClient{}
	}
	VaultUrl := env.VaultUrl
	SecretId := env.VaultSecretId
	RoleId := env.VaultRoleId
	logger.Info(fmt.Println("try to create Vault client"))
	client, err := vault.New(VaultUrl, SecretId, RoleId, env.VaultMountPoint)
	vaultClient := &VaultClient{
		Logger:    logger,
		LibClient: client,
	}
	if err != nil {
		vaultClient.LogErr(err)
	} else {
		logger.Info(fmt.Println("Vault client has been successfully created"))
	}
	return vaultClient
}

func (v VaultClient) LogErr(err error) {
	vaultConnection := fmt.Sprintf(
		"url:%s,mountpoint:%s,base_path:%s,role_id_id:%s",
		env.VaultUrl,
		env.VaultMountPoint,
		env.VaultSecretsBasePath,
		env.VaultRoleId,
	)
	v.Logger.Error(fmt.Sprintf("Valut error connection. Connection data: %s. Details: %s", vaultConnection, err))
}

// Получение секретов из vault
func (v VaultClient) GetVaultSecrets() (*VaultCredentials, error) {
	SecretStructure := VaultCredentials{}
	if v.LibClient.Vault == nil {
		return &SecretStructure, nil
	}
	data, err := v.LibClient.ReadSecretAsBytes(env.VaultSecretsBasePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &SecretStructure)
	if err != nil {
		return nil, err
	}
	return &SecretStructure, nil
}
