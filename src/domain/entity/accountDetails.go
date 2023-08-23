package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type AccountDetails struct {
	Username  valueObject.Username  `json:"username"`
	AccountId valueObject.AccountId `json:"id"`
	GroupId   valueObject.GroupId   `json:"groupId"`
}

func NewAccountDetails(
	username valueObject.Username,
	accountId valueObject.AccountId,
	groupId valueObject.GroupId,
) AccountDetails {
	return AccountDetails{
		Username:  username,
		AccountId: accountId,
		GroupId:   groupId,
	}
}
