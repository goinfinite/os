package useCase

import (
	"log"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

const (
	SslValidationsPerHour int    = 3
	OwnershipValidatePath string = "/validateOwnership"
)

type SslCertificateWatchDog struct {
	sslQueryRepo   repository.SslQueryRepo
	sslCmdRepo     repository.SslCmdRepo
	vhostQueryRepo repository.VirtualHostQueryRepo
	vhostCmdRepo   repository.VirtualHostCmdRepo
}

func NewSslCertificateWatchDog(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
) SslCertificateWatchDog {
	return SslCertificateWatchDog{
		sslQueryRepo:   sslQueryRepo,
		sslCmdRepo:     sslCmdRepo,
		vhostQueryRepo: vhostQueryRepo,
		vhostCmdRepo:   vhostCmdRepo,
	}
}

func (uc SslCertificateWatchDog) createInlineHtmlMapping(
	vhost valueObject.Fqdn,
	ownershipHash string,
) error {
	path, _ := valueObject.NewMappingPath(OwnershipValidatePath)
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := valueObject.NewHttpResponseCode(200)
	inlineHtmlContent, _ := valueObject.NewInlineHtmlContent(ownershipHash)

	inlineHmtlMapping := dto.NewCreateMapping(
		vhost,
		path,
		matchPattern,
		targetType,
		nil,
		nil,
		&httpResponseCode,
		&inlineHtmlContent,
	)
	return uc.vhostCmdRepo.CreateMapping(inlineHmtlMapping)
}

func (uc SslCertificateWatchDog) Execute() {
	vhosts, err := uc.vhostQueryRepo.Get()
	if err != nil {
		log.Printf("FailedToGetVhosts: %s", err.Error())
		return
	}

	invalidSslVhosts := []valueObject.Fqdn{}
	for _, vhost := range vhosts {
		isSslValid := uc.sslQueryRepo.IsSslPairValid(vhost.Hostname)
		if isSslValid {
			continue
		}

		invalidSslVhosts = append(invalidSslVhosts, vhost.Hostname)
	}

	for _, invalidSslVhost := range invalidSslVhosts {
		sslPair, err := uc.sslQueryRepo.GetSslPairByHostname(invalidSslVhost)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", invalidSslVhost, err.Error())
		}

		ownershipHash := uc.sslQueryRepo.GetOwnershipHash(sslPair.Certificate.CertificateContent)

		err = uc.createInlineHtmlMapping(invalidSslVhost, ownershipHash)
		if err != nil {
			log.Printf("FailedToCreateOwnershipValidationMapping: %s", err.Error())
			continue
		}

		// Wait for NGINX reload
		time.Sleep(5 * time.Second)

		isOwnershipValid := uc.sslQueryRepo.ValidateSslOwnership(invalidSslVhost, ownershipHash)
		if !isOwnershipValid {
			log.Printf("HostIsNotDomainOwner: %s", invalidSslVhost.String())
		}

		vhostMappings, err := uc.vhostQueryRepo.GetMappingsByHostname(invalidSslVhost)
		if err != nil {
			log.Printf("FailedToGetVhostMappings: %s", err.Error())
		}

		lastMappingIndex := len(vhostMappings) - 1
		lastMapping := vhostMappings[lastMappingIndex]
		err = uc.vhostCmdRepo.DeleteMapping(lastMapping)
		if err != nil {
			log.Printf("FailedToDeleteOwnershipValidationMapping: %s", err.Error())
		}

		if !isOwnershipValid {
			continue
		}

		err = uc.sslCmdRepo.Delete(sslPair.Id)
		if err != nil {
			log.Printf("FailedToDeleteInvalidSsl: %s", err.Error())
			continue
		}
	}
}
