package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
)

type Account struct {
	ID                     uint64 `gorm:"primaryKey"`
	GroupId                uint64 `gorm:"not null"`
	Username               string `gorm:"not null"`
	KeyHash                *string
	HomeDirectory          string `gorm:"not null"`
	IsSuperAdmin           bool   `gorm:"not null"`
	SecureAccessPublicKeys []SecureAccessPublicKey
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) ToModel(entity entity.Account) (model Account, err error) {
	return Account{
		ID:            entity.Id.Uint64(),
		GroupId:       entity.GroupId.Uint64(),
		Username:      entity.Username.String(),
		KeyHash:       nil,
		HomeDirectory: entity.HomeDirectory.String(),
		IsSuperAdmin:  entity.IsSuperAdmin,
	}, nil
}

func (model Account) ToEntity() (accountEntity entity.Account, err error) {
	accountId, err := valueObject.NewAccountId(model.ID)
	if err != nil {
		return accountEntity, err
	}

	groupId, err := valueObject.NewGroupId(model.GroupId)
	if err != nil {
		return accountEntity, err
	}

	username, err := valueObject.NewUsername(model.Username)
	if err != nil {
		return accountEntity, err
	}

	homeDirectory, err := valueObject.NewUnixFilePath(
		infraEnvs.UserDataBaseDirectory + "/" + username.String(),
	)
	if err != nil {
		return accountEntity, err
	}

	secureAccessPublicKeys := []entity.SecureAccessPublicKey{}
	for _, secureAccessPublicKeyModel := range model.SecureAccessPublicKeys {
		secureAccessPUblicKeyEntity, err := secureAccessPublicKeyModel.ToEntity()
		if err != nil {
			return accountEntity, err
		}
		secureAccessPublicKeys = append(
			secureAccessPublicKeys, secureAccessPUblicKeyEntity,
		)
	}

	return entity.NewAccount(
		accountId, groupId, username, homeDirectory, model.IsSuperAdmin,
		secureAccessPublicKeys, valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
