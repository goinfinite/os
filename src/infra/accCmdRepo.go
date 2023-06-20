package infra

import (
	"errors"
	"log"
	"os/exec"

	"github.com/speedianet/sam/src/domain/dto"
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
