package secureAccessKeyInfra

import (
	"errors"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type SecureAccessKeyCmdRepo struct {
	persistentDbSvc          *internalDbInfra.PersistentDatabaseService
	secureAccessKeyQueryRepo *SecureAccessKeyQueryRepo
	accountQueryRepo         *accountInfra.AccountQueryRepo
}

func NewSecureAccessKeyCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SecureAccessKeyCmdRepo {
	return &SecureAccessKeyCmdRepo{
		persistentDbSvc:          persistentDbSvc,
		secureAccessKeyQueryRepo: NewSecureAccessKeyQueryRepo(persistentDbSvc),
		accountQueryRepo:         accountInfra.NewAccountQueryRepo(persistentDbSvc),
	}
}

func (repo *SecureAccessKeyCmdRepo) createSecureAccessKeysFileIfNotExists(
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

func (repo *SecureAccessKeyCmdRepo) recreateSecureAccessKeysFile(
	accountId valueObject.AccountId,
	accountUsername valueObject.Username,
) error {
	keysFilePath := "/home/" + accountUsername.String() + "/.ssh/authorized_keys"
	if !infraHelper.FileExists(keysFilePath) {
		err := repo.createSecureAccessKeysFileIfNotExists(accountUsername)
		if err != nil {
			return err
		}
	}

	readRequestDto := dto.ReadSecureAccessKeysRequest{
		Pagination: dto.Pagination{
			ItemsPerPage: 1000,
		},
		AccountId: accountId,
	}
	keys, err := repo.secureAccessKeyQueryRepo.Read(readRequestDto)
	if err != nil {
		return err
	}

	keysFileContent := ""
	for _, key := range keys.SecureAccessKeys {
		keysFileContent += key.Content.String() + " " + key.Name.String() + "\n"
	}

	shouldOverwrite := true
	err = infraHelper.UpdateFile(keysFilePath, keysFileContent, shouldOverwrite)
	if err != nil {
		return errors.New("UpdateSecureAccessKeysFileContentError: " + err.Error())
	}

	return nil
}

func (repo *SecureAccessKeyCmdRepo) Create(
	createDto dto.CreateSecureAccessKey,
) (keyId valueObject.SecureAccessKeyId, err error) {
	accountEntity, err := repo.accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		return keyId, errors.New("AccountNotFound")
	}

	err = repo.createSecureAccessKeysFileIfNotExists(accountEntity.Username)
	if err != nil {
		return keyId, err
	}

	keyContentStr := createDto.Content.String()
	rawFingerprint, err := infraHelper.RunCmdWithSubShell(
		"echo \"" + keyContentStr + "\" | ssh-keygen -lf /dev/stdin | awk '{print $2}'",
	)
	if err != nil {
		return keyId, errors.New("ReadSecureAccessKeyFingerprintError: " + err.Error())
	}
	fingerPrint, err := valueObject.NewSecureAccessKeyFingerprint(rawFingerprint)
	if err != nil {
		return keyId, err
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"echo \"" + keyContentStr + "\" >> /home/" + accountEntity.Username.String() +
			"/.ssh/authorized_keys",
	)
	if err != nil {
		return keyId, errors.New("FailToAddNewSecureAccessKeyToFile: " + err.Error())
	}

	err = repo.allowAccountSecureRemoteConnection(accountEntity.Id)
	if err != nil {
		return keyId, err
	}

	secureAccessKeyModel := dbModel.NewSecureAccessKey(
		0, accountEntity.Id.Uint64(), createDto.Name.String(),
		createDto.Content.ReadWithoutKeyName(), fingerPrint.String(),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&secureAccessKeyModel)
	if createResult.Error != nil {
		return keyId, createResult.Error
	}

	keyId, err = valueObject.NewSecureAccessKeyId(secureAccessKeyModel.ID)
	if err != nil {
		return keyId, err
	}

	return keyId, repo.recreateSecureAccessKeysFile(
		accountEntity.Id, accountEntity.Username,
	)
}

func (repo *SecureAccessKeyCmdRepo) Delete(
	secureAccessKeyId valueObject.SecureAccessKeyId,
) error {
	readFirstRequestDto := dto.ReadSecureAccessKeysRequest{
		SecureAccessKeyId: &secureAccessKeyId,
	}
	keyToDelete, err := repo.secureAccessKeyQueryRepo.ReadFirst(readFirstRequestDto)
	if err != nil {
		return errors.New("SecureAccessKeyNotFound")
	}

	account, err := repo.accountQueryRepo.ReadById(keyToDelete.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = repo.persistentDbSvc.Handler.Delete(
		dbModel.SecureAccessKey{}, secureAccessKeyId.Uint16(),
	).Error
	if err != nil {
		return err
	}

	return repo.recreateSecureAccessKeysFile(account.Id, account.Username)
}
