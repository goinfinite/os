package accountInfra

import (
	"errors"
	"log/slog"
	"os"
	"os/user"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	accountQueryRepo *AccountQueryRepo
}

func NewAccountCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *AccountCmdRepo {
	return &AccountCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		accountQueryRepo: NewAccountQueryRepo(persistentDbSvc),
	}
}

func (repo *AccountCmdRepo) Create(
	createDto dto.CreateAccount,
) (accountId valueObject.AccountId, err error) {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(createDto.Password.String()), bcrypt.DefaultCost,
	)
	if err != nil {
		return accountId, errors.New("PasswordHashError: " + err.Error())
	}

	usernameStr := createDto.Username.String()

	_, err = infraHelper.RunCmd(
		"useradd", "-m",
		"-s", "/usr/sbin/nologin",
		"-p", string(passHash),
		usernameStr,
	)
	if err != nil {
		return accountId, errors.New("UserAddFailed: " + err.Error())
	}

	userInfo, err := user.Lookup(usernameStr)
	if err != nil {
		return accountId, errors.New("UserLookupFailed: " + err.Error())
	}

	accountId, err = valueObject.NewAccountId(userInfo.Uid)
	if err != nil {
		return accountId, err
	}

	groupId, err := valueObject.NewGroupId(userInfo.Gid)
	if err != nil {
		return accountId, err
	}

	nowUnixTime := valueObject.NewUnixTimeNow()
	accountEntity := entity.NewAccount(
		accountId, groupId, createDto.Username, nowUnixTime, nowUnixTime,
	)

	accountModel, err := dbModel.Account{}.ToModel(accountEntity)
	if err != nil {
		return accountId, err
	}

	err = repo.persistentDbSvc.Handler.Create(&accountModel).Error
	if err != nil {
		return accountId, err
	}

	return accountId, nil
}

func (repo *AccountCmdRepo) readUsernameById(
	accountId valueObject.AccountId,
) (username valueObject.Username, err error) {
	accountEntity, err := repo.accountQueryRepo.ReadById(accountId)
	if err != nil {
		return username, err
	}

	return accountEntity.Username, nil
}

func (repo *AccountCmdRepo) Delete(accountId valueObject.AccountId) error {
	username, err := repo.readUsernameById(accountId)
	if err != nil {
		return err
	}

	accountIdStr := accountId.String()

	_, err = infraHelper.RunCmd("pgrep", "-u", accountIdStr)
	if err == nil {
		_, _ = infraHelper.RunCmd("pkill", "-9", "-U", accountIdStr)
	}

	_, err = infraHelper.RunCmd("userdel", "-r", username.String())
	if err != nil {
		return err
	}

	accountModel := dbModel.Account{}

	err = repo.persistentDbSvc.Handler.Delete(accountModel, accountIdStr).Error
	if err != nil {
		return errors.New("DeleteAccountDatabaseError: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) UpdatePassword(
	accountId valueObject.AccountId, password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()), bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("PasswordHashError: " + err.Error())
	}

	username, err := repo.readUsernameById(accountId)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("usermod", "-p", string(passHash), username.String())
	if err != nil {
		return errors.New("UserModFailed: " + err.Error())
	}

	accountModel := dbModel.Account{ID: accountId.Uint64()}
	return repo.persistentDbSvc.Handler.
		Model(&accountModel).
		Update("updated_at", time.Now()).
		Error
}

func (repo *AccountCmdRepo) UpdateApiKey(
	accountId valueObject.AccountId,
) (tokenValue valueObject.AccessTokenStr, err error) {
	uuidStr := uuid.New().String()
	apiKeyPlainText := accountId.String() + ":" + uuidStr

	secretKey := os.Getenv("ACCOUNT_API_KEY_SECRET")
	encryptedApiKey, err := infraHelper.EncryptStr(secretKey, apiKeyPlainText)
	if err != nil {
		return tokenValue, err
	}

	apiKey, err := valueObject.NewAccessTokenStr(encryptedApiKey)
	if err != nil {
		return tokenValue, err
	}

	uuidHash := infraHelper.GenStrongHash(uuidStr)

	accountModel := dbModel.Account{ID: accountId.Uint64()}
	updateResult := repo.persistentDbSvc.Handler.
		Model(&accountModel).
		Update("key_hash", uuidHash)
	if updateResult.Error != nil {
		return tokenValue, err
	}

	return apiKey, nil
}

func (repo *AccountCmdRepo) ensureSecureAccessKeysDirAndFileExistence(
	accountUsername valueObject.Username,
) error {
	accountUsernameStr := accountUsername.String()

	secureAccessKeysDirPath := "/home/" + accountUsernameStr + "/.ssh"
	if !infraHelper.FileExists(secureAccessKeysDirPath) {
		err := infraHelper.MakeDir(secureAccessKeysDirPath)
		if err != nil {
			return errors.New("CreateSecureAccessKeysDirectoryError: " + err.Error())
		}
	}

	secureAccessKeysFilePath := secureAccessKeysDirPath + "/authorized_keys"
	if infraHelper.FileExists(secureAccessKeysFilePath) {
		return nil
	}

	_, err := os.Create(secureAccessKeysFilePath)
	if err != nil {
		return errors.New("CreateSecureAccessKeysFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown", "-R", accountUsernameStr, secureAccessKeysFilePath,
	)
	if err != nil {
		return errors.New("ChownSecureAccessKeysFileError: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) isSecureAccessKeyValid(
	keyContent valueObject.SecureAccessKeyContent,
) bool {
	keyName, err := keyContent.ReadOnlyKeyName()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	keyNameStr := keyName.String()

	keyTempFilePath := "/tmp/" + keyNameStr + "_secureAccessKey"
	shouldOverwrite := true
	err = infraHelper.UpdateFile(
		keyTempFilePath, keyContent.String(), shouldOverwrite,
	)
	if err != nil {
		slog.Error(
			"CreateSecureAccessKeyTempFileError", slog.String("keyName", keyNameStr),
			slog.Any("err", err),
		)
		return false
	}

	_, err = infraHelper.RunCmdWithSubShell("ssh-keygen -l -f " + keyTempFilePath)
	if err != nil {
		slog.Error(
			"ValidateSecureAccessKeyError", slog.String("keyName", keyNameStr),
			slog.Any("err", err),
		)
		return false
	}

	err = os.Remove(keyTempFilePath)
	if err != nil {
		slog.Error(
			"DeleteSecureAccessKeyTempFileError", slog.String("keyName", keyNameStr),
			slog.Any("err", err),
		)
	}

	return true
}

func (repo *AccountCmdRepo) CreateSecureAccessKey(
	createDto dto.CreateSecureAccessKey,
) (keyId valueObject.SecureAccessKeyId, err error) {
	account, err := repo.accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		return keyId, errors.New("AccountNotFound")
	}

	err = repo.ensureSecureAccessKeysDirAndFileExistence(account.Username)
	if err != nil {
		return keyId, err
	}

	keyContentStr := createDto.Content.ReadWithoutKeyName() + " " +
		createDto.Name.String()
	keyContent, err := valueObject.NewSecureAccessKeyContent(keyContentStr)
	if err != nil {
		return keyId, errors.New("InvalidSecureAccessKey")
	}

	if !repo.isSecureAccessKeyValid(keyContent) {
		return keyId, errors.New("InvalidSecureAccessKey")
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"echo \"" + keyContentStr + "\" >> /home/" + account.Username.String() +
			"/.ssh/authorized_keys",
	)
	if err != nil {
		return keyId, errors.New("FailToAddNewSecureAccessKeyToFile: " + err.Error())
	}

	key, err := repo.accountQueryRepo.ReadSecureAccessKeyByName(
		createDto.AccountId, createDto.Name,
	)
	if err != nil {
		return keyId, errors.New("SecureAccessKeyWasNotCreated")
	}

	return key.Id, nil
}

func (repo *AccountCmdRepo) DeleteSecureAccessKey(
	deleteDto dto.DeleteSecureAccessKey,
) error {
	account, err := repo.accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	keyToDelete, err := repo.accountQueryRepo.ReadSecureAccessKeyById(
		deleteDto.AccountId, deleteDto.Id,
	)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"sed -i '\\|" + keyToDelete.Content.String() + "|d' " +
			"/home/" + account.Username.String() + "/.ssh/authorized_keys",
	)
	if err != nil {
		return errors.New("FailToDeleteSecureAccessKeyFromFile: " + err.Error())
	}

	return nil
}
