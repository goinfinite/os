package infra

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"golang.org/x/net/publicsuffix"
)

var configurationsDir string = "/app/conf/nginx"

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
		log.Println("InvalidServerName: " + serverNamesParts[0])
		return vhosts, nil
	}

	primaryDomain, _ := valueObject.NewFqdn(os.Getenv("VIRTUAL_HOST"))
	isPrimaryDomain := firstDomain == primaryDomain

	for _, serverName := range serverNamesParts {
		serverName, err := valueObject.NewFqdn(serverName)
		if err != nil {
			log.Println("InvalidServerName: " + serverName.String())
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

		rootDomainStr, err := publicsuffix.EffectiveTLDPlusOne(serverName.String())
		if err != nil {
			log.Println("InvalidRootDomain: " + serverName.String())
			continue
		}
		rootDomain, err := valueObject.NewFqdn(rootDomainStr)
		if err != nil {
			log.Println("InvalidRootDomain: " + rootDomainStr)
			continue
		}

		isSubdomain := rootDomain != serverName
		if isSubdomain {
			vhostType, _ = valueObject.NewVirtualHostType("subdomain")
			parentDomainPtr = &rootDomain
		}

		rootDirectorySuffix := "/" + serverName.String()
		if isPrimaryDomain {
			rootDirectorySuffix = ""
		}
		rootDirectory, err := valueObject.NewUnixFilePath(
			"/app/html" + rootDirectorySuffix,
		)
		if err != nil {
			log.Println("InvalidRootDirectory: " + rootDirectorySuffix)
			continue
		}

		if isAliases {
			vhostType, _ = valueObject.NewVirtualHostType("alias")
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

	configsDirHandler, err := os.Open(configurationsDir)
	if err != nil {
		log.Fatal(err)
	}
	defer configsDirHandler.Close()

	files, err := configsDirHandler.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".conf") {
			continue
		}
		filePath, err := valueObject.NewUnixFilePath(
			configurationsDir + "/" + fileName,
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
