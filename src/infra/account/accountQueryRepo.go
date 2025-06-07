package accountInfra

import (
	"errors"
	"log/slog"
	"math"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
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

func (repo *AccountQueryRepo) Count(
	requestDto dto.ReadAccountsRequest,
) (count uint64, err error) {
	model := dbModel.Account{}
	if requestDto.AccountId != nil {
		model.ID = requestDto.AccountId.Uint64()
	}
	if requestDto.AccountUsername != nil {
		model.Username = requestDto.AccountUsername.String()
	}

	var itemsTotal int64
	err = repo.persistentDbSvc.Handler.Model(&model).Where(&model).Count(&itemsTotal).Error
	if err != nil {
		return count, errors.New("CountAccountsTotalError: " + err.Error())
	}

	return uint64(itemsTotal), nil
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
	if requestDto.ShouldIncludeSecureAccessPublicKeys != nil && *requestDto.ShouldIncludeSecureAccessPublicKeys {
		dbQuery = dbQuery.Preload("SecureAccessPublicKeys")
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

	accountModels := []dbModel.Account{}
	err = dbQuery.Find(&accountModels).Error
	if err != nil {
		return responseDto, errors.New("ReadAccountsError: " + err.Error())
	}

	accountEntities := []entity.Account{}
	for _, model := range accountModels {
		accountEntity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"AccountModelToEntityError", slog.Uint64("id", uint64(model.ID)),
				slog.String("err", err.Error()),
			)
			continue
		}

		accountEntities = append(accountEntities, accountEntity)
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
		Accounts:   accountEntities,
	}, nil
}

func (repo *AccountQueryRepo) ReadFirst(
	requestDto dto.ReadAccountsRequest,
) (account entity.Account, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return account, err
	}

	if len(responseDto.Accounts) == 0 {
		return account, errors.New("AccountNotFound")
	}

	return responseDto.Accounts[0], nil
}

func (repo *AccountQueryRepo) ReadSecureAccessPublicKeys(
	requestDto dto.ReadSecureAccessPublicKeysRequest,
) (responseDto dto.ReadSecureAccessPublicKeysResponse, err error) {
	model := dbModel.SecureAccessPublicKey{
		AccountID: requestDto.AccountId.Uint64(),
	}
	if requestDto.SecureAccessPublicKeyId != nil {
		model.ID = requestDto.SecureAccessPublicKeyId.Uint16()
	}
	if requestDto.SecureAccessPublicKeyName != nil {
		model.Name = requestDto.SecureAccessPublicKeyName.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&model).
		Where(&model)

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New(
			"CountSecureAccessPublicKeysTotalError: " + err.Error(),
		)
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

	secureAccessKeysModels := []dbModel.SecureAccessPublicKey{}
	err = dbQuery.Find(&secureAccessKeysModels).Error
	if err != nil {
		return responseDto, errors.New("ReadSecureAccessPublicKeysError: " + err.Error())
	}

	secureAccessKeysEntities := []entity.SecureAccessPublicKey{}
	for _, model := range secureAccessKeysModels {
		secureAccessKeysEntity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"SecureAccessPublicKeyModelToEntityError",
				slog.Uint64("id", uint64(model.ID)), slog.String("err", err.Error()),
			)
			continue
		}

		secureAccessKeysEntities = append(
			secureAccessKeysEntities, secureAccessKeysEntity,
		)
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

	return dto.ReadSecureAccessPublicKeysResponse{
		Pagination:             responsePagination,
		SecureAccessPublicKeys: secureAccessKeysEntities,
	}, nil
}

func (repo *AccountQueryRepo) ReadFirstSecureAccessPublicKey(
	requestDto dto.ReadSecureAccessPublicKeysRequest,
) (secureAccessPublicKey entity.SecureAccessPublicKey, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadSecureAccessPublicKeys(requestDto)
	if err != nil {
		return secureAccessPublicKey, err
	}

	if len(responseDto.SecureAccessPublicKeys) == 0 {
		return secureAccessPublicKey, errors.New("SecureAccessPublicKeyNotFound")
	}

	return responseDto.SecureAccessPublicKeys[0], nil
}
