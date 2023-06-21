package infra

import (
	"errors"
	"log"

	"github.com/msteinert/pam"
	"github.com/speedianet/sam/src/domain/dto"
)

type AuthQueryRepo struct {
}

func (repo AuthQueryRepo) IsLoginValid(login dto.Login) bool {
	tx, err := pam.StartFunc(
		"system-auth",
		login.Username.String(),
		func(s pam.Style, msg string) (string, error) {
			switch s {
			case pam.PromptEchoOff:
				return login.Password.String(), nil
			case pam.PromptEchoOn, pam.ErrorMsg, pam.TextInfo:
				log.Printf("PamMessage: %s", msg)
				return "", nil
			}
			log.Println("UnhandledPamMessageStyle")
			return "", errors.New("UnhandledPamMessageStyle")
		})

	if err != nil {
		log.Printf("PamError: %v", err)
		return false
	}

	if err = tx.Authenticate(0); err != nil {
		return false
	}

	return true
}
