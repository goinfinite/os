package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessPublicKey struct {
	ID        uint16 `gorm:"primarykey"`
	AccountId uint64 `gorm:"not null"`
	Name      string `gorm:"not null"`
	Content   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SecureAccessPublicKey) TableName() string {
	return "secure_access_public_keys"
}

func NewSecureAccessPublicKey(
	id uint16,
	accountId uint64,
	name, content string,
) SecureAccessPublicKey {
	model := SecureAccessPublicKey{
		AccountId: accountId,
		Name:      name,
		Content:   content,
	}

	if id != 0 {
		model.ID = id
	}

	return model
}

func (model SecureAccessPublicKey) ToEntity() (
	SecureAccessPublicKeyEntity entity.SecureAccessPublicKey, err error,
) {
	id, err := valueObject.NewSecureAccessPublicKeyId(model.ID)
	if err != nil {
		return SecureAccessPublicKeyEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountId)
	if err != nil {
		return SecureAccessPublicKeyEntity, err
	}

	name, err := valueObject.NewSecureAccessPublicKeyName(model.Name)
	if err != nil {
		return SecureAccessPublicKeyEntity, err
	}

	content, err := valueObject.NewSecureAccessPublicKeyContent(model.Content)
	if err != nil {
		return SecureAccessPublicKeyEntity, err
	}

	return entity.NewSecureAccessPublicKey(
		id, accountId, content, &name,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	)
}
