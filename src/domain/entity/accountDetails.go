package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type AccountDetails struct {
	Username valueObject.Username `json:"username"`
	UserId   valueObject.UserId   `json:"id"`
	GroupId  valueObject.GroupId  `json:"groupId"`
}

func NewAccountDetails(
	username valueObject.Username,
	userId valueObject.UserId,
	groupId valueObject.GroupId,
) AccountDetails {
	return AccountDetails{
		Username: username,
		UserId:   userId,
		GroupId:  groupId,
	}
}
