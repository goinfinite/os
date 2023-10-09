package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddSsl(
	sslCmdRepo repository.SslCmdRepo,
	addSsl dto.AddSsl,
) error {
	err := sslCmdRepo.Add(addSsl)
	if err != nil {
		log.Printf("AddSslError: %s", err)
		return errors.New("AddSslInfraError")
	}

	log.Printf(
		"SSL added to '%v' cname in '%v' virtual host.",
		addSsl.Certificate.CommonName,
		addSsl.Hostname.String(),
	)

	return nil
}
