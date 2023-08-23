package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          valueObject.AccountId `json:"accountId"`
	Password           *valueObject.Password `json:"password"`
	ShouldUpdateApiKey *bool                 `json:"shouldUpdateApiKey"`
}

func NewUpdateAccount(
	accountId valueObject.AccountId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
	}
}
