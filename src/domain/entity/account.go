package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type Account struct {
	Id       valueObject.AccountId `json:"id"`
	GroupId  valueObject.GroupId   `json:"groupId"`
	Username valueObject.Username  `json:"username"`
}

func NewAccount(
	username valueObject.Username,
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
) Account {
	return Account{
		Id:       accountId,
		GroupId:  groupId,
		Username: username,
	}
}
