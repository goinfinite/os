package accountInfra

import (
	"errors"
	"os"
	"os/user"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
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

func (repo *AccountCmdRepo) createAuthorizedKeysFile(
	accountUsername valueObject.Username,
	accountHomeDirectory valueObject.UnixFilePath,
) error {
	accountUsernameStr := accountUsername.String()

	sshDirPath := accountHomeDirectory.String() + "/.ssh"
	err := infraHelper.MakeDir(sshDirPath)
	if err != nil {
		return errors.New("CreateSshDirectoryError: " + err.Error())
	}

	authorizedKeysFilePath := sshDirPath + "/authorized_keys"
	_, err = os.Create(authorizedKeysFilePath)
	if err != nil {
		return errors.New("CreateAuthorizedKeysFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown", "-R", accountUsernameStr, authorizedKeysFilePath,
	)
	if err != nil {
		return errors.New("ChownAuthorizedKeysFileError: " + err.Error())
	}

	return nil
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
	homeDirectory, err := valueObject.NewUnixFilePath(
		infraEnvs.UserDataBaseDirectory + "/" + usernameStr,
	)
	if err != nil {
		return accountId, errors.New("DefineHomeDirectoryError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"useradd", "-m",
		"-s", "/bin/bash",
		"-p", string(passHash),
		usernameStr,
	)
	if err != nil {
		return accountId, errors.New("UserAddFailed: " + err.Error())
	}

	err = repo.createAuthorizedKeysFile(createDto.Username, homeDirectory)
	if err != nil {
		return accountId, err
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
		accountId, groupId, createDto.Username, homeDirectory,
		[]entity.SecureAccessPublicKey{}, nowUnixTime, nowUnixTime,
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

func (repo *AccountCmdRepo) Delete(accountId valueObject.AccountId) error {
	readFirstAccountRequestDto := dto.ReadAccountsRequest{
		AccountId: &accountId,
	}
	accountEntity, err := repo.accountQueryRepo.ReadFirst(readFirstAccountRequestDto)
	if err != nil {
		return err
	}

	accountIdStr := accountId.String()

	_, err = infraHelper.RunCmd("pgrep", "-u", accountIdStr)
	if err == nil {
		_, _ = infraHelper.RunCmd("pkill", "-9", "-U", accountIdStr)
	}

	_, err = infraHelper.RunCmd("userdel", "-r", accountEntity.Username.String())
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

	readFirstAccountRequestDto := dto.ReadAccountsRequest{
		AccountId: &accountId,
	}
	accountEntity, err := repo.accountQueryRepo.ReadFirst(readFirstAccountRequestDto)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(
		"usermod", "-p", string(passHash), accountEntity.Username.String(),
	)
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

func (repo *AccountCmdRepo) rebuildAuthorizedKeysFile(
	accountId valueObject.AccountId,
	accountHomeDirectory valueObject.UnixFilePath,
) error {
	readPublicKeysRequestDto := dto.ReadSecureAccessPublicKeysRequest{
		Pagination: dto.Pagination{
			ItemsPerPage: 1000,
		},
		AccountId: accountId,
	}
	readPublicKeysResponseDto, err := repo.accountQueryRepo.ReadSecureAccessPublicKeys(
		readPublicKeysRequestDto,
	)
	if err != nil {
		return err
	}

	keysFileContent := "# Please, don't edit manually as this will be automatically recreated.\n\n"
	for _, keyEntity := range readPublicKeysResponseDto.SecureAccessPublicKeys {
		keysFileContent += keyEntity.Content.String() + " " +
			keyEntity.Name.String() + "\n"
	}

	authorizedKeysFilePath := accountHomeDirectory.String() + "/.ssh/authorized_keys"
	shouldOverwrite := true
	err = infraHelper.UpdateFile(
		authorizedKeysFilePath, keysFileContent, shouldOverwrite,
	)
	if err != nil {
		return errors.New("UpdateAuthorizedKeysFileContentError: " + err.Error())
	}

	return nil
}

func (repo *AccountCmdRepo) CreateSecureAccessPublicKey(
	createDto dto.CreateSecureAccessPublicKey,
) (keyId valueObject.SecureAccessPublicKeyId, err error) {
	readFirstAccountRequestDto := dto.ReadAccountsRequest{
		AccountId: &createDto.AccountId,
	}
	accountEntity, err := repo.accountQueryRepo.ReadFirst(readFirstAccountRequestDto)
	if err != nil {
		return keyId, errors.New("AccountNotFound")
	}

	secureAccessPublicKeyModel := dbModel.NewSecureAccessPublicKey(
		0, accountEntity.Id.Uint64(), createDto.Name.String(),
		createDto.Content.ReadWithoutKeyName(),
	)

	dbCreateResult := repo.persistentDbSvc.Handler.Create(&secureAccessPublicKeyModel)
	if dbCreateResult.Error != nil {
		return keyId, dbCreateResult.Error
	}

	keyId, err = valueObject.NewSecureAccessPublicKeyId(secureAccessPublicKeyModel.ID)
	if err != nil {
		return keyId, err
	}

	return keyId, repo.rebuildAuthorizedKeysFile(
		accountEntity.Id, accountEntity.HomeDirectory,
	)
}

func (repo *AccountCmdRepo) DeleteSecureAccessPublicKey(
	secureAccessPublicKeyId valueObject.SecureAccessPublicKeyId,
) error {
	readFirstPublicKeyRequestDto := dto.ReadSecureAccessPublicKeysRequest{
		SecureAccessPublicKeyId: &secureAccessPublicKeyId,
	}
	keyEntity, err := repo.accountQueryRepo.ReadFirstSecureAccessPublicKey(
		readFirstPublicKeyRequestDto,
	)
	if err != nil {
		return errors.New("SecureAccessPublicKeyNotFound")
	}

	readFirstAccountRequestDto := dto.ReadAccountsRequest{
		AccountId: &keyEntity.AccountId,
	}
	accountEntity, err := repo.accountQueryRepo.ReadFirst(readFirstAccountRequestDto)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = repo.persistentDbSvc.Handler.Delete(
		dbModel.SecureAccessPublicKey{}, secureAccessPublicKeyId.Uint16(),
	).Error
	if err != nil {
		return err
	}

	return repo.rebuildAuthorizedKeysFile(
		accountEntity.Id, accountEntity.HomeDirectory,
	)
}
