package infra

import (
	"errors"
	"log"
	"os/exec"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	"golang.org/x/crypto/bcrypt"
)

type AccCmdRepo struct {
}

func (repo AccCmdRepo) Add(addUser dto.AddUser) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(addUser.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("PasswordHashError: %s", err)
		return errors.New("PasswordHashError")
	}

	addUserCmd := exec.Command(
		"useradd",
		"-m",
		"-s", "/bin/bash",
		"-p", string(passHash),
		addUser.Username.String(),
	)

	err = addUserCmd.Run()
	if err != nil {
		log.Printf("UserAddError: %s", err)
		return errors.New("UserAddError")
	}

	return nil
}

func (repo AccCmdRepo) Delete(userId valueObject.UserId) error {
	accQuery := AccQueryRepo{}
	accDetails, err := accQuery.GetById(userId)
	if err != nil {
		log.Printf("GetUserDetailsError: %s", err)
		return errors.New("GetUserDetailsError")
	}

	delUserCmd := exec.Command(
		"userdel",
		"-r",
		accDetails.Username.String(),
	)

	err = delUserCmd.Run()
	if err != nil {
		log.Printf("UserDeleteError: %s", err)
		return errors.New("UserDeleteError")
	}

	return nil
}
