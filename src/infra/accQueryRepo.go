package infra

import (
	"errors"
	"os/user"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AccQueryRepo struct {
}

func accDetailsFactory(userInfo *user.User) (entity.AccountDetails, error) {
	username, err := valueObject.NewUsername(userInfo.Username)
	if err != nil {
		return entity.AccountDetails{}, errors.New("UsernameParseError")
	}

	userId, err := valueObject.NewUserIdFromString(userInfo.Uid)
	if err != nil {
		return entity.AccountDetails{}, errors.New("UserIdParseError")
	}

	groupId, err := valueObject.NewGroupIdFromString(userInfo.Gid)
	if err != nil {
		return entity.AccountDetails{}, errors.New("GroupIdParseError")
	}

	return entity.NewAccountDetails(
		username,
		userId,
		groupId,
	), nil
}

func (repo AccQueryRepo) GetByUsername(
	username valueObject.Username,
) (entity.AccountDetails, error) {
	userInfo, err := user.Lookup(string(username))
	if err != nil {
		return entity.AccountDetails{}, errors.New("UserLookupError")
	}

	return accDetailsFactory(userInfo)
}

func (repo AccQueryRepo) GetById(
	userId valueObject.UserId,
) (entity.AccountDetails, error) {
	userInfo, err := user.LookupId(userId.String())
	if err != nil {
		return entity.AccountDetails{}, errors.New("UserLookupError")
	}

	return accDetailsFactory(userInfo)
}
