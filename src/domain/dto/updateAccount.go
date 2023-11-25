package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateAccount struct {
	Id                 valueObject.AccountId `json:"id"`
	Password           *valueObject.Password `json:"password"`
	ShouldUpdateApiKey *bool                 `json:"shouldUpdateApiKey"`
}

func NewUpdateAccount(
	id valueObject.AccountId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
) UpdateAccount {
	return UpdateAccount{
		Id:                 id,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
	}
}
