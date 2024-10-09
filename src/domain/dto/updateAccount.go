package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateAccount struct {
	Id                 *valueObject.AccountId `json:"id"`
	Username           *valueObject.Username  `json:"username"`
	Password           *valueObject.Password  `json:"password"`
	ShouldUpdateApiKey *bool                  `json:"shouldUpdateApiKey"`
}

func NewUpdateAccount(
	id *valueObject.AccountId,
	username *valueObject.Username,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
) UpdateAccount {
	return UpdateAccount{
		Id:                 id,
		Username:           username,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
	}
}
