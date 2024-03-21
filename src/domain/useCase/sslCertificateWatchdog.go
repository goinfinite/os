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

type SslCertificateWatchdog struct {
	sslQueryRepo   repository.SslQueryRepo
	sslCmdRepo     repository.SslCmdRepo
	vhostQueryRepo repository.VirtualHostQueryRepo
	vhostCmdRepo   repository.VirtualHostCmdRepo
}

func NewSslCertificateWatchdog(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
) SslCertificateWatchdog {
	return SslCertificateWatchdog{
		sslQueryRepo:   sslQueryRepo,
		sslCmdRepo:     sslCmdRepo,
		vhostQueryRepo: vhostQueryRepo,
		vhostCmdRepo:   vhostCmdRepo,
	}
}

func (uc SslCertificateWatchdog) createInlineHtmlMapping(
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

func (uc SslCertificateWatchdog) Execute() {
	sslPairs, err := uc.sslQueryRepo.GetSslPairs()
	if err != nil {
		log.Printf("FailedToGetSslPairs: %s", err.Error())
		return
	}

	for _, sslPair := range sslPairs {
		isSslValid := uc.sslQueryRepo.IsSslPairValid(sslPair)
		if isSslValid {
			continue
		}

		ownershipHash := uc.sslQueryRepo.GetOwnershipHash(sslPair.Certificate.CertificateContent)

		firstVhost := sslPair.VirtualHosts[0]
		err = uc.createInlineHtmlMapping(firstVhost, ownershipHash)
		if err != nil {
			log.Printf("FailedToCreateOwnershipValidationMapping: %s", err.Error())
			continue
		}

		// Wait for NGINX reload
		time.Sleep(5 * time.Second)

		isOwnershipValid := uc.vhostQueryRepo.IsDomainOwner(firstVhost, ownershipHash)
		if !isOwnershipValid {
			log.Printf("CurrentHostIsNotDomainOwner: %s", firstVhost.String())
		}

		vhostMappings, err := uc.vhostQueryRepo.GetMappingsByHostname(firstVhost)
		if err != nil {
			log.Printf("FailedToGetVhostMappings: %s", err.Error())
			continue
		}

		if len(vhostMappings) == 0 {
			log.Printf("VhostMappingsNotFound: %s", firstVhost)
			continue
		}

		lastMappingIndex := len(vhostMappings) - 1
		lastMapping := vhostMappings[lastMappingIndex]
		err = uc.vhostCmdRepo.DeleteMapping(lastMapping)
		if err != nil {
			log.Printf("FailedToDeleteOwnershipValidationMapping: %s", err.Error())
			continue
		}

		if !isOwnershipValid {
			continue
		}

		err = uc.sslCmdRepo.ReplaceWithValidSsl(firstVhost)
		if err != nil {
			log.Printf("%s: %s", firstVhost.String(), err.Error())
		}
	}
}
