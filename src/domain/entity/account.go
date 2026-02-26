package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type Account struct {
	Id                     tkValueObject.AccountId            `json:"id"`
	GroupId                tkValueObject.UnixGroupId          `json:"groupId"`
	Username               valueObject.Username               `json:"username"`
	HomeDirectory          tkValueObject.UnixAbsoluteFilePath `json:"homeDirectory"`
	IsSuperAdmin           bool                               `json:"isSuperAdmin"`
	SecureAccessPublicKeys []SecureAccessPublicKey            `json:"secureAccessPublicKeys"`
	CreatedAt              tkValueObject.UnixTime             `json:"createdAt"`
	UpdatedAt              tkValueObject.UnixTime             `json:"updatedAt"`
}

func NewAccount(
	accountId tkValueObject.AccountId,
	groupId tkValueObject.UnixGroupId,
	username valueObject.Username,
	homeDirectory tkValueObject.UnixAbsoluteFilePath,
	isSuperAdmin bool,
	secureAccessPublicKeys []SecureAccessPublicKey,
	createdAt, updatedAt tkValueObject.UnixTime,
) Account {
	return Account{
		Id:                     accountId,
		GroupId:                groupId,
		Username:               username,
		HomeDirectory:          homeDirectory,
		IsSuperAdmin:           isSuperAdmin,
		SecureAccessPublicKeys: secureAccessPublicKeys,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}
