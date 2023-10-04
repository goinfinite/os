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

	sslCertShortVersion := addSsl.Certificate.String()[:75]
	log.Printf("SSL '%v' added to '%v' virtual host.", sslCertShortVersion, addSsl.Hostname.String())

	return nil
}
