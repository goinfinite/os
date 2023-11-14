package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Account struct {
	Id       valueObject.AccountId `json:"id"`
	GroupId  valueObject.GroupId   `json:"groupId"`
	Username valueObject.Username  `json:"username"`
}

func NewAccount(
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
	username valueObject.Username,
) Account {
	return Account{
		Id:       accountId,
		GroupId:  groupId,
		Username: username,
	}
}
