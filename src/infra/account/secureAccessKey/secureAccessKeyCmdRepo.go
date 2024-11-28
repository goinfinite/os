package secureAccessKeyInfra

import (
	"errors"
	"log/slog"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

type SecureAccessKeyCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	accountQueryRepo *accountInfra.AccountQueryRepo
}

func NewSecureAccessKeyCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SecureAccessKeyCmdRepo {
	return &SecureAccessKeyCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		accountQueryRepo: accountInfra.NewAccountQueryRepo(persistentDbSvc),
	}
}

func (repo *SecureAccessKeyCmdRepo) ensureSecureAccessKeysDirAndFileExistence(
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

func (repo *SecureAccessKeyCmdRepo) isSecureAccessKeyValid(
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

func (repo *SecureAccessKeyCmdRepo) allowAccountSecureRemoteConnection(
	accountId valueObject.AccountId,
) error {
	accountUsername, err := infraHelper.RunCmdWithSubShell(
		"awk -F: '$3 == " + accountId.String() +
			" && $7 != \"/bin/bash\" {print $1}' /etc/passwd",
	)
	if err != nil {
		return errors.New("ReadUnixUsernameFromFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"chsh -s /bin/bash " + accountUsername,
	)
	if err != nil {
		return errors.New("ChangeDefaultBashError: " + err.Error())
	}

	return nil
}

func (repo *SecureAccessKeyCmdRepo) Create(
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

	err = repo.allowAccountSecureRemoteConnection(createDto.AccountId)
	if err != nil {
		return keyId, err
	}

	return keyId, nil
}

func (repo *SecureAccessKeyCmdRepo) Delete(
	deleteDto dto.DeleteSecureAccessKey,
) error {
	_, err := repo.accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	/*keyToDelete, err := repo.secureAccessKeyQueryRepo.ReadById(
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
	}*/

	return nil
}
