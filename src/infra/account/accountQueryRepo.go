package accountInfra

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type AccountQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewAccountQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *AccountQueryRepo {
	return &AccountQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *AccountQueryRepo) Read() ([]entity.Account, error) {
	accountEntities := []entity.Account{}

	var accountModels []dbModel.Account
	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.Account{}).
		Find(&accountModels).Error
	if err != nil {
		return accountEntities, errors.New("QueryAccountsError: " + err.Error())
	}

	for _, accountModel := range accountModels {
		accountEntity, err := accountModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Any("error", err.Error()),
				slog.Uint64("accountId", accountModel.ID),
			)
			continue
		}

		accountEntities = append(accountEntities, accountEntity)
	}

	return accountEntities, nil
}

func (repo *AccountQueryRepo) ReadById(
	accountId valueObject.AccountId,
) (accountEntity entity.Account, err error) {
	var accountModel dbModel.Account
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.Account{}).
		Where("id = ?", accountId.String()).
		Find(&accountModel).Error
	if err != nil {
		return accountEntity, errors.New("QueryAccountByIdError: " + err.Error())
	}

	return accountModel.ToEntity()
}

func (repo *AccountQueryRepo) ReadByUsername(
	accountUsername valueObject.Username,
) (accountEntity entity.Account, err error) {
	var accountModel dbModel.Account
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.Account{}).
		Where("username = ?", accountUsername.String()).
		Find(&accountModel).Error
	if err != nil {
		return accountEntity, errors.New("QueryAccountByUsernameError: " + err.Error())
	}

	return accountModel.ToEntity()
}

func (repo *AccountQueryRepo) secureAccessKeyFactory(
	rawSecureAccessKeyContent string,
	secureAccessKeySecret string,
) (secureAccessKey entity.SecureAccessKey, err error) {
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

	rawKeyHashId, err := infraHelper.EncryptStr(
		secureAccessKeySecret, keyContent.String(),
	)
	if err != nil {
		return secureAccessKey, err
	}

	keyHashId, err := valueObject.NewSecureAccessKeyHashId(rawKeyHashId)
	if err != nil {
		return secureAccessKey, err
	}

	return entity.NewSecureAccessKey(keyHashId, keyName, keyContent), nil
}

func (repo *AccountQueryRepo) ReadSecureAccessKeys(
	accountId valueObject.AccountId,
) ([]entity.SecureAccessKey, error) {
	secureAccessKeys := []entity.SecureAccessKey{}

	account, err := repo.ReadById(accountId)
	if err != nil {
		return secureAccessKeys, errors.New("AccountNotFound")
	}
	accountUsernameStr := account.Username.String()

	secureAccessKeysDirPath := "/home/" + accountUsernameStr + "/.ssh"
	err = infraHelper.MakeDir(secureAccessKeysDirPath)
	if err != nil {
		return secureAccessKeys, errors.New(
			"CreateSecureAccessKeysDirectoryError: " + err.Error(),
		)
	}

	secureAccessKeysFilePath := secureAccessKeysDirPath + "/authorized_keys"
	if !infraHelper.FileExists(secureAccessKeysFilePath) {
		_, err = os.Create(secureAccessKeysFilePath)
		if err != nil {
			return secureAccessKeys, errors.New(
				"CreateSecureAccessKeysFileError: " + err.Error(),
			)
		}

		_, err = infraHelper.RunCmd(
			"chown", "-R", accountUsernameStr, secureAccessKeysFilePath,
		)
		if err != nil {
			return secureAccessKeys, errors.New(
				"ChownSecureAccessKeysFileError: " + err.Error(),
			)
		}
	}

	secureAccessKeysFileContent, err := infraHelper.GetFileContent(
		secureAccessKeysFilePath,
	)
	if err != nil {
		return secureAccessKeys, errors.New(
			"ReadSecureAccessKeysFileContentError: " + err.Error(),
		)
	}

	secretKey := os.Getenv("ACCOUNT_SECURE_ACCESS_KEY_SECRET")

	secureAccessKeysFileContentParts := strings.Split(secureAccessKeysFileContent, "\n")
	for index, rawSecureAccessKeyContent := range secureAccessKeysFileContentParts {
		secureAccessKey, err := repo.secureAccessKeyFactory(
			rawSecureAccessKeyContent, secretKey,
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", index))
			continue
		}

		secureAccessKeys = append(secureAccessKeys, secureAccessKey)
	}

	return secureAccessKeys, nil
}
