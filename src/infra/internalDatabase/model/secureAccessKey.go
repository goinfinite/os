package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessKey struct {
	ID        uint16 `gorm:"primarykey"`
	AccountId uint64 `gorm:"not null"`
	Name      string `gorm:"primarykey"`
	Content   string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SecureAccessKey) TableName() string {
	return "secure_access_key"
}

func NewSecureAccessKey(
	id uint16,
	accountId uint64,
	name, content string,
) SecureAccessKey {
	model := SecureAccessKey{
		AccountId: accountId,
		Name:      name,
		Content:   content,
	}

	if id != 0 {
		model.ID = id
	}

	return model
}

func (model SecureAccessKey) ToEntity() (
	secureAccessKeyEntity entity.SecureAccessKey, err error,
) {
	id, err := valueObject.NewSecureAccessKeyId(model.ID)
	if err != nil {
		return secureAccessKeyEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountId)
	if err != nil {
		return secureAccessKeyEntity, err
	}

	name, err := valueObject.NewSecureAccessKeyName(model.Name)
	if err != nil {
		return secureAccessKeyEntity, err
	}

	content, err := valueObject.NewSecureAccessKeyContent(model.Content)
	if err != nil {
		return secureAccessKeyEntity, err
	}

	return entity.NewSecureAccessKey(
		id, accountId, name, content,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
