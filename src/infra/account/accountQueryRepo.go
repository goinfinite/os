package accountInfra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
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
