package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type Account struct {
	Username  valueObject.Username  `json:"username"`
	AccountId valueObject.AccountId `json:"id"`
	GroupId   valueObject.GroupId   `json:"groupId"`
}

func NewAccount(
	username valueObject.Username,
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
) Account {
	return Account{
		Username:  username,
		AccountId: accountId,
		GroupId:   groupId,
	}
}
