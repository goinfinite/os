package secureAccessKeyInfra

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

type SecureAccessKeyQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewSecureAccessKeyQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SecureAccessKeyQueryRepo {
	return &SecureAccessKeyQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *SecureAccessKeyQueryRepo) Read(
	requestDto dto.ReadSecureAccessPublicKeysRequest,
) (responseDto dto.ReadSecureAccessPublicKeysResponse, err error) {
	model := dbModel.SecureAccessPublicKey{
		AccountId: requestDto.AccountId.Uint64(),
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

	models := []dbModel.SecureAccessPublicKey{}
	err = dbQuery.Find(&models).Error
	if err != nil {
		return responseDto, errors.New("ReadSecureAccessPublicKeysError: " + err.Error())
	}

	entities := []entity.SecureAccessPublicKey{}
	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"SecureAccessPublicKeyModelToEntityError",
				slog.Uint64("id", uint64(model.ID)), slog.Any("error", err),
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

	return dto.ReadSecureAccessPublicKeysResponse{
		Pagination:             responsePagination,
		SecureAccessPublicKeys: entities,
	}, nil
}

func (repo *SecureAccessKeyQueryRepo) ReadFirst(
	requestDto dto.ReadSecureAccessPublicKeysRequest,
) (secureAccessPublicKey entity.SecureAccessPublicKey, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return secureAccessPublicKey, err
	}

	if len(responseDto.SecureAccessPublicKeys) == 0 {
		return secureAccessPublicKey, errors.New("SecureAccessPublicKeyNotFound")
	}

	return responseDto.SecureAccessPublicKeys[0], nil
}
