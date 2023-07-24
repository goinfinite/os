package infra

import (
	"errors"
	"os/user"
	"strings"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
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

func (repo AccQueryRepo) Get() ([]entity.AccountDetails, error) {
	output, err := infraHelper.RunCmd("awk", "-F:", "{print $1}", "/etc/passwd")
	if err != nil {
		return []entity.AccountDetails{}, errors.New("UsersLookupError")
	}

	usernames := strings.Split(string(output), "\n")
	var accsDetails []entity.AccountDetails
	for _, username := range usernames {
		username, err := valueObject.NewUsername(username)
		if err != nil {
			continue
		}

		accDetails, err := repo.GetByUsername(username)
		if err != nil {
			continue
		}
		if accDetails.UserId < 1000 {
			continue
		}
		if accDetails.Username == "nobody" {
			continue
		}

		accsDetails = append(accsDetails, accDetails)
	}

	return accsDetails, nil
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
