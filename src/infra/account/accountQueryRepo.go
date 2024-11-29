package accountInfra

import (
	"errors"
	"log/slog"
	"math"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"github.com/iancoleman/strcase"
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

func (repo *AccountQueryRepo) Read(
	requestDto dto.ReadAccountsRequest,
) (responseDto dto.ReadAccountsResponse, err error) {
	model := dbModel.Account{}
	if requestDto.AccountId != nil {
		model.ID = requestDto.AccountId.Uint64()
	}
	if requestDto.AccountUsername != nil {
		model.Username = requestDto.AccountUsername.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&model).
		Where(&model)
	if requestDto.ShouldIncludeSecureAccessKeys != nil && *requestDto.ShouldIncludeSecureAccessKeys {
		dbQuery = dbQuery.Preload("SecureAccessKeys")
	}

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New("CountAccountsTotalError: " + err.Error())
	}

	dbQuery.Limit(int(requestDto.Pagination.ItemsPerPage))
	if requestDto.Pagination.LastSeenId == nil {
		offset := int(requestDto.Pagination.PageNumber) * int(requestDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", requestDto.Pagination.LastSeenId.String())
	}
	if requestDto.Pagination.SortBy != nil {
		orderStatement := requestDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if requestDto.Pagination.SortDirection != nil {
			orderStatement += " " + requestDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	models := []dbModel.Account{}
	err = dbQuery.Find(&models).Error
	if err != nil {
		return responseDto, errors.New("ReadAccountsError: " + err.Error())
	}

	entities := []entity.Account{}
	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"AccountModelToEntityError", slog.Uint64("id", uint64(model.ID)),
				slog.Any("error", err),
			)
			continue
		}

		entities = append(entities, entity)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(requestDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    requestDto.Pagination.PageNumber,
		ItemsPerPage:  requestDto.Pagination.ItemsPerPage,
		SortBy:        requestDto.Pagination.SortBy,
		SortDirection: requestDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadAccountsResponse{
		Pagination: responsePagination,
		Accounts:   entities,
	}, nil
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
