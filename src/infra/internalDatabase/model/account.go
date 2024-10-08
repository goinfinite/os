package dbModel

import (
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type Account struct {
	ID        uint64 `gorm:"primarykey"`
	GroupId   uint64 `gorm:"not null"`
	Username  string `gorm:"not null"`
	KeyHash   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) ToModel(entity entity.Account) (model Account, err error) {
	return Account{
		ID:       entity.Id.Uint64(),
		GroupId:  entity.GroupId.Uint64(),
		Username: entity.Username.String(),
		KeyHash:  nil,
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

	return entity.NewAccount(
		accountId, groupId, username,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
