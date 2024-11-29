package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type Account struct {
	Id               valueObject.AccountId `json:"id"`
	GroupId          valueObject.GroupId   `json:"groupId"`
	Username         valueObject.Username  `json:"username"`
	SecureAccessKeys []SecureAccessKey     `json:"secureAccessKeys"`
	CreatedAt        valueObject.UnixTime  `json:"createdAt"`
	UpdatedAt        valueObject.UnixTime  `json:"updatedAt"`
}

func NewAccount(
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
	username valueObject.Username,
	secureAccessKeys []SecureAccessKey,
	createdAt, updatedAt valueObject.UnixTime,
) Account {
	return Account{
		Id:               accountId,
		GroupId:          groupId,
		Username:         username,
		SecureAccessKeys: secureAccessKeys,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
