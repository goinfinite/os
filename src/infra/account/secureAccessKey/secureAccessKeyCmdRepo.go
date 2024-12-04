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

func (repo *SecureAccessKeyCmdRepo) createSecureAccessPublicKeysFileIfNotExists(
	accountUsername valueObject.Username,
) error {
	accountUsernameStr := accountUsername.String()

	secureAccessPublicKeysDirPath := "/home/" + accountUsernameStr + "/.ssh"
	if !infraHelper.FileExists(secureAccessPublicKeysDirPath) {
		err := infraHelper.MakeDir(secureAccessPublicKeysDirPath)
		if err != nil {
			return errors.New("CreateSecureAccessPublicKeysDirectoryError: " + err.Error())
		}
	}

	secureAccessPublicKeysFilePath := secureAccessPublicKeysDirPath + "/authorized_keys"
	if infraHelper.FileExists(secureAccessPublicKeysFilePath) {
		return nil
	}

	_, err := os.Create(secureAccessPublicKeysFilePath)
	if err != nil {
		return errors.New("CreateSecureAccessPublicKeysFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown", "-R", accountUsernameStr, secureAccessPublicKeysFilePath,
	)
	if err != nil {
		return errors.New("ChownSecureAccessPublicKeysFileError: " + err.Error())
	}

	return nil
}

func (repo *SecureAccessKeyCmdRepo) recreateSecureAccessPublicKeysFile(
	accountId valueObject.AccountId,
	accountUsername valueObject.Username,
) error {
	keysFilePath := "/home/" + accountUsername.String() + "/.ssh/authorized_keys"
	if !infraHelper.FileExists(keysFilePath) {
		err := repo.createSecureAccessPublicKeysFileIfNotExists(accountUsername)
		if err != nil {
			return err
		}
	}

	readRequestDto := dto.ReadSecureAccessPublicKeysRequest{
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
	for _, key := range keys.SecureAccessPublicKeys {
		keysFileContent += key.Content.String() + " " + key.Name.String() + "\n"
	}

	shouldOverwrite := true
	err = infraHelper.UpdateFile(keysFilePath, keysFileContent, shouldOverwrite)
	if err != nil {
		return errors.New("UpdateSecureAccessPublicKeysFileContentError: " + err.Error())
	}

	return nil
}

func (repo *SecureAccessKeyCmdRepo) Create(
	createDto dto.CreateSecureAccessPublicKey,
) (keyId valueObject.SecureAccessPublicKeyId, err error) {
	accountEntity, err := repo.accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		return keyId, errors.New("AccountNotFound")
	}

	err = repo.createSecureAccessPublicKeysFileIfNotExists(accountEntity.Username)
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
	fingerPrint, err := valueObject.NewSecureAccessPublicKeyFingerprint(rawFingerprint)
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

	secureAccessKeyModel := dbModel.NewSecureAccessPublicKey(
		0, accountEntity.Id.Uint64(), createDto.Name.String(),
		createDto.Content.ReadWithoutKeyName(), fingerPrint.String(),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&secureAccessKeyModel)
	if createResult.Error != nil {
		return keyId, createResult.Error
	}

	keyId, err = valueObject.NewSecureAccessPublicKeyId(secureAccessKeyModel.ID)
	if err != nil {
		return keyId, err
	}

	return keyId, repo.recreateSecureAccessPublicKeysFile(
		accountEntity.Id, accountEntity.Username,
	)
}

func (repo *SecureAccessKeyCmdRepo) Delete(
	secureAccessPublicKeyId valueObject.SecureAccessPublicKeyId,
) error {
	readFirstRequestDto := dto.ReadSecureAccessPublicKeysRequest{
		SecureAccessPublicKeyId: &secureAccessPublicKeyId,
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
		dbModel.SecureAccessPublicKey{}, secureAccessPublicKeyId.Uint16(),
	).Error
	if err != nil {
		return err
	}

	return repo.recreateSecureAccessPublicKeysFile(account.Id, account.Username)
}
