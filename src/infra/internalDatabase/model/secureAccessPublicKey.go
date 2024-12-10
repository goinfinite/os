package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessPublicKey struct {
	ID        uint16 `gorm:"primarykey"`
	AccountID uint64 `gorm:"not null"`
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
		AccountID: accountId,
		Name:      name,
		Content:   content,
	}

	if id != 0 {
		model.ID = id
	}

	return model
}

func (model SecureAccessPublicKey) ToEntity() (
	secureAccessPublicKeyEntity entity.SecureAccessPublicKey, err error,
) {
	id, err := valueObject.NewSecureAccessPublicKeyId(model.ID)
	if err != nil {
		return secureAccessPublicKeyEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return secureAccessPublicKeyEntity, err
	}

	content, err := valueObject.NewSecureAccessPublicKeyContent(model.Content)
	if err != nil {
		return secureAccessPublicKeyEntity, err
	}

	fingerprint, err := content.ReadFingerprint()
	if err != nil {
		return secureAccessPublicKeyEntity, err
	}

	name, err := valueObject.NewSecureAccessPublicKeyName(model.Name)
	if err != nil {
		return secureAccessPublicKeyEntity, err
	}

	return entity.NewSecureAccessPublicKey(
		id, accountId, content, fingerprint, name,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
