package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type Account struct {
	Id        valueObject.AccountId `json:"id"`
	GroupId   valueObject.GroupId   `json:"groupId"`
	Username  valueObject.Username  `json:"username"`
	CreatedAt valueObject.UnixTime  `json:"createdAt"`
	UpdatedAt valueObject.UnixTime  `json:"updatedAt"`
}

func NewAccount(
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
	username valueObject.Username,
	createdAt, updatedAt valueObject.UnixTime,
) Account {
	return Account{
		Id:        accountId,
		GroupId:   groupId,
		Username:  username,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
