package useCase

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func getDefaultStartupFileByMultiService(
	serviceName valueObject.ServiceName,
) (valueObject.UnixFilePath, error) {
	switch serviceName.String() {
	case "node":
		return valueObject.NewUnixFilePath("/app/html/index.js")
	default:
		return "", errors.New("UnknownInstallableMultiService")
	}
}

func getServiceNameWithSuffix(
	startupFile valueObject.UnixFilePath,
	serviceName valueObject.ServiceName,
) (valueObject.ServiceName, error) {
	startupFileBytes := []byte(startupFile.String())
	startupFileHash := md5.Sum(startupFileBytes)
	startupFileHashStr := hex.EncodeToString(startupFileHash[:])
	startupFileShortHashStr := startupFileHashStr[:12]

	svcNameWithSuffix := serviceName.String() + "-" + startupFileShortHashStr
	return valueObject.NewServiceName(svcNameWithSuffix)
}

func AddInstallableService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addDto dto.AddInstallableService,
) error {
	_, err := servicesQueryRepo.GetByName(addDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installableSvcs, err := servicesQueryRepo.GetInstallables()
	if err != nil {
		log.Printf("GetInstallableServicesError: %s", err.Error())
		return errors.New("GetInstallableServicesInfraError")
	}

	isNatureMulti := false
	for _, installableSvc := range installableSvcs {
		if installableSvc.Name.String() == addDto.Name.String() {
			isNatureMulti = installableSvc.Nature.String() == "multi"
		}
	}

	if isNatureMulti {
		startupFile, err := getDefaultStartupFileByMultiService(addDto.Name)
		if err != nil {
			return err
		}

		if addDto.StartupFile != nil {
			startupFile = *addDto.StartupFile
		}

		newSvcName, err := getServiceNameWithSuffix(startupFile, addDto.Name)
		if err != nil {
			return err
		}

		addDto.Name = newSvcName
	}

	err = servicesCmdRepo.AddInstallable(addDto)
	if err != nil {
		log.Printf("AddInstallableServiceError: %s", err.Error())
		return errors.New("AddInstallableServiceInfraError")
	}

	vhostsWithMappings, err := vhostQueryRepo.GetWithMappings()
	if err != nil {
		log.Printf("GetVhostsWithMappingError: %s", err.Error())
		return errors.New("GetVhostsWithMappingsInfraError")
	}

	if len(vhostsWithMappings) == 0 {
		return errors.New("VhostsNotFound")
	}

	primaryVhostWithMapping := vhostsWithMappings[0]
	shouldCreateFirstMapping := len(primaryVhostWithMapping.Mappings) == 0 && addDto.AutoCreateMapping
	if !shouldCreateFirstMapping {
		return nil
	}

	serviceMapping, err := serviceMappingFactory(
		primaryVhostWithMapping.Hostname,
		addDto.Name,
	)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingError")
	}

	err = vhostCmdRepo.AddMapping(serviceMapping)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingInfraError")
	}

	return nil
}
