package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type Account struct {
	Id                     valueObject.AccountId   `json:"id"`
	GroupId                valueObject.GroupId     `json:"groupId"`
	Username               valueObject.Username    `json:"username"`
	SecureAccessPublicKeys []SecureAccessPublicKey `json:"secureAccessPublicKeys"`
	CreatedAt              valueObject.UnixTime    `json:"createdAt"`
	UpdatedAt              valueObject.UnixTime    `json:"updatedAt"`
}

func NewAccount(
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
	username valueObject.Username,
	secureAccessPublicKeys []SecureAccessPublicKey,
	createdAt, updatedAt valueObject.UnixTime,
) Account {
	return Account{
		Id:                     accountId,
		GroupId:                groupId,
		Username:               username,
		SecureAccessPublicKeys: secureAccessPublicKeys,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}
}
