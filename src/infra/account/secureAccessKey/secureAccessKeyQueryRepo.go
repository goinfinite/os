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
	requestDto dto.ReadSecureAccessKeysRequest,
) (responseDto dto.ReadSecureAccessKeysResponse, err error) {
	model := dbModel.SecureAccessKey{
		AccountId: requestDto.AccountId.Uint64(),
	}
	if requestDto.SecureAccessKeyId != nil {
		model.ID = requestDto.SecureAccessKeyId.Uint16()
	}
	if requestDto.SecureAccessKeyName != nil {
		model.Name = requestDto.SecureAccessKeyName.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&model).
		Where(&model)

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New(
			"CountSecureAccessKeysTotalError: " + err.Error(),
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

	models := []dbModel.SecureAccessKey{}
	err = dbQuery.Find(&models).Error
	if err != nil {
		return responseDto, errors.New("ReadSecureAccessKeysError: " + err.Error())
	}

	entities := []entity.SecureAccessKey{}
	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"SecureAccessKeyModelToEntityError",
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

	return dto.ReadSecureAccessKeysResponse{
		Pagination:       responsePagination,
		SecureAccessKeys: entities,
	}, nil
}

func (repo *SecureAccessKeyQueryRepo) ReadFirst(
	requestDto dto.ReadSecureAccessKeysRequest,
) (secureAccessKey entity.SecureAccessKey, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return secureAccessKey, err
	}

	if len(responseDto.SecureAccessKeys) == 0 {
		return secureAccessKey, errors.New("SecureAccessKeyNotFound")
	}

	return responseDto.SecureAccessKeys[0], nil
}
