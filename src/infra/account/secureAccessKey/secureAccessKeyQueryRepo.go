package secureAccessKeyInfra

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

type SecureAccessKeyQueryRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	accountQueryRepo *accountInfra.AccountQueryRepo
}

func NewSecureAccessKeyQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SecureAccessKeyQueryRepo {
	return &SecureAccessKeyQueryRepo{
		persistentDbSvc:  persistentDbSvc,
		accountQueryRepo: accountInfra.NewAccountQueryRepo(persistentDbSvc),
	}
}

func (repo *SecureAccessKeyQueryRepo) secureAccessKeyFactory(
	rawKeyId int,
	rawSecureAccessKeyContent string,
	accountId valueObject.AccountId,
) (secureAccessKey entity.SecureAccessKey, err error) {
	keyId, err := valueObject.NewSecureAccessKeyId(rawKeyId)
	if err != nil {
		return secureAccessKey, err
	}

	keyContent, err := valueObject.NewSecureAccessKeyContent(
		rawSecureAccessKeyContent,
	)
	if err != nil {
		return secureAccessKey, err
	}

	keyName, err := keyContent.ReadOnlyKeyName()
	if err != nil {
		return secureAccessKey, err
	}

	now := valueObject.NewUnixTimeNow()

	return entity.NewSecureAccessKey(
		keyId, accountId, keyName, keyContent, now, now,
	), nil
}

func (repo *SecureAccessKeyQueryRepo) Read(
	accountId valueObject.AccountId,
) ([]entity.SecureAccessKey, error) {
	secureAccessKeys := []entity.SecureAccessKey{}

	account, err := repo.accountQueryRepo.ReadById(accountId)
	if err != nil {
		return secureAccessKeys, errors.New("AccountNotFound")
	}

	secureAccessKeysFilePath := "/home/" + account.Username.String() + "/.ssh" +
		"/authorized_keys"
	secureAccessKeysFileContent, err := infraHelper.GetFileContent(
		secureAccessKeysFilePath,
	)
	if err != nil {
		return secureAccessKeys, errors.New(
			"ReadSecureAccessKeysFileContentError: " + err.Error(),
		)
	}

	secureAccessKeysFileContentParts := strings.Split(secureAccessKeysFileContent, "\n")
	for index, rawSecureAccessKeyContent := range secureAccessKeysFileContentParts {
		if rawSecureAccessKeyContent == "" {
			continue
		}

		rawKeyId := index + 1
		secureAccessKey, err := repo.secureAccessKeyFactory(
			rawKeyId, rawSecureAccessKeyContent, accountId,
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", index))
			continue
		}

		secureAccessKeys = append(secureAccessKeys, secureAccessKey)
	}

	return secureAccessKeys, nil
}

func (repo *SecureAccessKeyQueryRepo) ReadById(
	accountId valueObject.AccountId,
	secureAccessKeyId valueObject.SecureAccessKeyId,
) (secureAccessKey entity.SecureAccessKey, err error) {
	secureAccessKeys, err := repo.Read(accountId)
	if err != nil {
		return secureAccessKey, err
	}

	for _, key := range secureAccessKeys {
		if key.Id.Uint16() != secureAccessKeyId.Uint16() {
			continue
		}

		return key, nil
	}

	return secureAccessKey, errors.New("SecureAccessKeyNotFound")
}

func (repo *SecureAccessKeyQueryRepo) ReadByName(
	accountId valueObject.AccountId,
	secureAccessKeyName valueObject.SecureAccessKeyName,
) (secureAccessKey entity.SecureAccessKey, err error) {
	secureAccessKeys, err := repo.Read(accountId)
	if err != nil {
		return secureAccessKey, err
	}

	for _, key := range secureAccessKeys {
		if key.Name.String() != secureAccessKeyName.String() {
			continue
		}

		return key, nil
	}

	return secureAccessKey, errors.New("SecureAccessKeyNotFound")
}
