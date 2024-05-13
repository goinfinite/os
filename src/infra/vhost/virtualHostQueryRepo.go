package vhostInfra

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	infraData "github.com/speedianet/os/src/infra/infraData"
)

type VirtualHostQueryRepo struct {
}

func (repo VirtualHostQueryRepo) vhostsFactory(
	filePath valueObject.UnixFilePath,
) ([]entity.VirtualHost, error) {
	vhosts := []entity.VirtualHost{}

	fileContent, err := infraHelper.GetFileContent(filePath.String())
	if err != nil {
		return vhosts, err
	}

	serverNamesRegex := regexp.MustCompile(`(?m)^\s*server_name\s+(.+);$`)
	serverNamesMatches := serverNamesRegex.FindStringSubmatch(fileContent)
	if len(serverNamesMatches) == 0 {
		return vhosts, errors.New("GetServerNameFailed")
	}

	serverNamesParts := strings.Split(serverNamesMatches[1], " ")
	if len(serverNamesParts) == 0 {
		return vhosts, errors.New("GetServerNameFailed")
	}

	firstDomain, err := valueObject.NewFqdn(serverNamesParts[0])
	if err != nil {
		log.Printf("InvalidServerName: %s", serverNamesParts[0])
		return vhosts, nil
	}
	isPrimaryVhost := infraHelper.IsPrimaryVirtualHost(firstDomain)

	for _, serverName := range serverNamesParts {
		serverName, err := valueObject.NewFqdn(serverName)
		if err != nil {
			log.Printf("InvalidServerName: %s", serverName.String())
			continue
		}

		isWww := strings.HasPrefix(serverName.String(), "www.")
		if isWww {
			continue
		}

		var parentDomainPtr *valueObject.Fqdn
		vhostType, _ := valueObject.NewVirtualHostType("top-level")
		isAliases := serverName != firstDomain
		if isAliases {
			vhostType, _ = valueObject.NewVirtualHostType("alias")
			parentDomainPtr = &firstDomain
		}

		rootDomain, err := infraHelper.GetRootDomain(serverName)
		if err != nil {
			log.Printf("%s: %s", err.Error(), serverName.String())
			continue
		}

		isSubdomain := rootDomain != serverName
		if isSubdomain {
			vhostType, _ = valueObject.NewVirtualHostType("subdomain")
			parentDomainPtr = &rootDomain
		}

		if isPrimaryVhost {
			vhostType, _ = valueObject.NewVirtualHostType("primary")
		}

		rootDirectorySuffix := "/" + serverName.String()
		if isPrimaryVhost {
			rootDirectorySuffix = ""
		}
		rootDirectory, err := valueObject.NewUnixFilePath(
			infraData.GlobalConfigs.PrimaryPublicDir + rootDirectorySuffix,
		)
		if err != nil {
			log.Printf("InvalidRootDirectory: %s", rootDirectorySuffix)
			continue
		}

		vhost := entity.NewVirtualHost(
			serverName,
			vhostType,
			rootDirectory,
			parentDomainPtr,
		)

		vhosts = append(vhosts, vhost)
	}

	return vhosts, nil
}

func (repo VirtualHostQueryRepo) Get() ([]entity.VirtualHost, error) {
	vhostsList := []entity.VirtualHost{}

	configsDirHandler, err := os.Open(infraData.GlobalConfigs.VirtualHostsConfDir)
	if err != nil {
		return vhostsList, errors.New("FailedToOpenConfDir: " + err.Error())
	}
	defer configsDirHandler.Close()

	files, err := configsDirHandler.Readdir(-1)
	if err != nil {
		return vhostsList, errors.New("FailedToReadConfDir: " + err.Error())
	}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".conf") {
			continue
		}
		filePath, err := valueObject.NewUnixFilePath(
			infraData.GlobalConfigs.VirtualHostsConfDir + "/" + fileName,
		)
		if err != nil {
			log.Println("InvalidVirtualHostFile: " + fileName)
			continue
		}

		vhosts, err := repo.vhostsFactory(filePath)
		if err != nil {
			log.Println("VirtualHostFileParseError: " + fileName)
			continue
		}
		vhostsList = append(vhostsList, vhosts...)
	}

	return vhostsList, nil
}

func (repo VirtualHostQueryRepo) GetByHostname(
	hostname valueObject.Fqdn,
) (entity.VirtualHost, error) {
	var virtualHost entity.VirtualHost

	vhosts, err := repo.Get()
	if err != nil {
		return virtualHost, err
	}

	for _, vhost := range vhosts {
		if vhost.Hostname == hostname {
			return vhost, nil
		}
	}

	return virtualHost, errors.New("VirtualHostNotFound")
}

func (repo VirtualHostQueryRepo) GetVirtualHostMappingsFilePath(
	vhostName valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var mappingFilePath valueObject.UnixFilePath

	mappingFileName := vhostName.String() + ".conf"

	vhostEntity, err := repo.GetByHostname(vhostName)
	if err != nil {
		return mappingFilePath, errors.New("VirtualHostNotFound")
	}

	isAlias := vhostEntity.Type.String() == "alias"
	if isAlias {
		parentHostname := *vhostEntity.ParentHostname
		mappingFileName = parentHostname.String() + ".conf"
	}

	if infraHelper.IsPrimaryVirtualHost(vhostName) {
		mappingFileName = "primary.conf"
	}

	return valueObject.NewUnixFilePath(
		infraData.GlobalConfigs.MappingsConfDir + "/" + mappingFileName,
	)
}
