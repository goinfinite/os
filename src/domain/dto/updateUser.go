package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type UpdateUser struct {
	UserId             valueObject.UserId    `json:"userId"`
	Password           *valueObject.Password `json:"password"`
	ShouldUpdateApiKey *bool                 `json:"shouldUpdateApiKey"`
}

func NewUpdateUser(
	userId valueObject.UserId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
) UpdateUser {
	return UpdateUser{
		UserId:             userId,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
	}
}
