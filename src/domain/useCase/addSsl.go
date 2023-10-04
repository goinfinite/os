package useCase

import (
	"errors"
	"log"
	"strings"

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
	parsedSslCertShortVersion := strings.Replace(sslCertShortVersion, "\n", "\\n", -1)
	log.Printf("SSL '%v' added to '%v' virtual host.", parsedSslCertShortVersion, addSsl.Hostname.String())

	return nil
}
