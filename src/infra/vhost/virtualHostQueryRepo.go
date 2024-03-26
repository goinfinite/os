package vhostInfra

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	envDataInfra "github.com/speedianet/os/src/infra/shared"
	"golang.org/x/exp/slices"
)

var configurationsDir string = "/app/conf/nginx"
var mappingsDir string = "/app/conf/nginx/mapping"

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
	isPrimaryDomain := infraHelper.IsPrimaryVirtualHost(firstDomain)

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

		if isPrimaryDomain {
			vhostType, _ = valueObject.NewVirtualHostType("primary")
		}

		rootDirectorySuffix := "/" + serverName.String()
		if isPrimaryDomain {
			rootDirectorySuffix = ""
		}
		rootDirectory, err := valueObject.NewUnixFilePath(
			"/app/html" + rootDirectorySuffix,
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

	configsDirHandler, err := os.Open(configurationsDir)
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

	return valueObject.NewUnixFilePath(mappingsDir + "/" + mappingFileName)
}

func (repo VirtualHostQueryRepo) locationBlockToMapping(
	locationBlockIndex int,
	locationBlockParts []string,
	vhostHost valueObject.Fqdn,
	serviceNamePortsMap map[string][]string,
) (entity.Mapping, error) {
	var mapping entity.Mapping

	if len(locationBlockParts) < 3 {
		return mapping, errors.New("GetLocationBlockPartsFailed")
	}

	modifierAndPath := locationBlockParts[1]
	modifierAndPathParts := strings.Split(modifierAndPath, " ")
	if len(modifierAndPathParts) == 0 {
		return mapping, errors.New("GetModifierAndPathPartsFailed")
	}

	modifier := ""
	pathStr := modifierAndPathParts[0]
	if len(modifierAndPathParts) == 2 {
		modifier = modifierAndPathParts[0]
		pathStr = modifierAndPathParts[1]
	}

	validModifiers := []string{"=", "~"}
	isModifierEmpty := modifier == ""
	isModifierValid := slices.Contains(validModifiers, modifier)

	if !isModifierEmpty && !isModifierValid {
		return mapping, errors.New("InvalidModifier: " + modifier)
	}

	matchPatternStr := "begins-with"
	isModifierEquals := modifier == "="
	if isModifierEquals {
		matchPatternStr = "equals"
	}

	isModifierRegex := modifier == "~"
	if isModifierRegex {
		matchPatternStr = "contains"

		lastPathCharIsDollarSign := strings.HasSuffix(pathStr, "$")
		if lastPathCharIsDollarSign {
			pathStr = strings.TrimSuffix(pathStr, "$")
			matchPatternStr = "ends-with"
		}
	}

	matchPattern, err := valueObject.NewMappingMatchPattern(matchPatternStr)
	if err != nil {
		return mapping, errors.New("InvalidMatchPattern: " + matchPatternStr)
	}

	path, err := valueObject.NewMappingPath(pathStr)
	if err != nil {
		return mapping, errors.New("InvalidMappingPath: " + pathStr)
	}

	targetTypeStr := "service"

	blockContent := locationBlockParts[2]
	blockContent = strings.TrimSpace(blockContent)

	isStaticFiles := strings.Contains(blockContent, "try_files")
	if isStaticFiles {
		targetTypeStr = "static-files"
	}

	var targetUrlPtr *valueObject.Url
	var targetResponseCodePtr *valueObject.HttpResponseCode
	var targetInlineHtmlContentPtr *valueObject.InlineHtmlContent

	isUrlOrResponseCodeOrInlineHtml := strings.Contains(blockContent, "return ")
	if isUrlOrResponseCodeOrInlineHtml {
		blockContentLines := strings.Split(blockContent, "\n")

		blockContentFirstLine := blockContentLines[0]
		directiveBlockContent := blockContentFirstLine
		directiveBlockContentParts := strings.Split(directiveBlockContent, " ")
		if len(directiveBlockContentParts) < 2 {
			return mapping, errors.New("GetLocationBlockContentPartsFailed")
		}

		targetTypeStr = "response-code"

		directive := directiveBlockContentParts[0]
		if directive == "add_header" {
			blockContentSecondLine := blockContentLines[1]
			blockContentSecondLine = strings.TrimSpace(blockContentSecondLine)
			blockContentSecondLineParts := strings.Split(blockContentSecondLine, " ")
			if len(blockContentSecondLineParts) < 2 {
				return mapping, errors.New("GetLocationBlockContentPartsFailed")
			}

			inlineHtmlContentStr := blockContentSecondLineParts[2]
			inlineHtmlContentWithoutQuotesStr := strings.ReplaceAll(inlineHtmlContentStr, "'", "")
			inlineHtmlContentWithoutSemicolonStr := strings.TrimRight(
				inlineHtmlContentWithoutQuotesStr, ";",
			)
			inlineHtmlContentStr = inlineHtmlContentWithoutSemicolonStr
			inlineHtmlContent, err := valueObject.NewInlineHtmlContent(inlineHtmlContentStr)
			if err != nil {
				return mapping, errors.New("InvalidReturnInlineHtmlContent: " + inlineHtmlContentStr)
			}
			targetInlineHtmlContentPtr = &inlineHtmlContent

			directive = blockContentSecondLineParts[0]
			directiveBlockContentParts = blockContentSecondLineParts
			targetTypeStr = "inline-html"
		}

		if directive != "return" {
			return mapping, errors.New("GetLocationDirectiveFailed")
		}

		responseCodeStr := directiveBlockContentParts[1]
		if targetTypeStr == "response-code" {
			responseCodeWithoutSemicolonStr := strings.TrimRight(responseCodeStr, ";")
			responseCodeStr = responseCodeWithoutSemicolonStr
		}

		if len(responseCodeStr) == 0 {
			return mapping, errors.New("InvalidReturnResponseCode: " + responseCodeStr)
		}

		responseCode, err := valueObject.NewHttpResponseCode(responseCodeStr)
		if err != nil {
			return mapping, errors.New("InvalidReturnResponseCode: " + responseCodeStr)
		}
		targetResponseCodePtr = &responseCode

		hasUrl := len(directiveBlockContentParts) == 3 && targetTypeStr != "inline-html"
		if hasUrl {
			targetTypeStr = "url"

			urlStr := directiveBlockContentParts[2]
			urlStr = strings.TrimSuffix(urlStr, ";")
			url, err := valueObject.NewUrl(urlStr)
			if err != nil {
				return mapping, errors.New("InvalidReturnUrl: " + urlStr)
			}
			targetUrlPtr = &url
		}
	}

	targetType, err := valueObject.NewMappingTargetType(targetTypeStr)
	if err != nil {
		return mapping, errors.New("InvalidTargetType: " + targetTypeStr)
	}

	var targetServiceNamePtr *valueObject.ServiceName
	if targetTypeStr == "service" {
		hostnamePortRegex := regexp.MustCompile(`(?m)localhost:\d{1,5}`)
		hostnamePortMatches := hostnamePortRegex.FindStringSubmatch(blockContent)
		if len(hostnamePortMatches) == 0 {
			return mapping, errors.New("GetServicePortsFailed")
		}

		for _, hostnamePortMatch := range hostnamePortMatches {
			hostnamePortParts := strings.Split(hostnamePortMatch, ":")
			if len(hostnamePortParts) != 2 {
				continue
			}

			port := hostnamePortParts[1]
			for serviceName, ports := range serviceNamePortsMap {
				if !slices.Contains(ports, port) {
					continue
				}
				serviceName, _ := valueObject.NewServiceName(serviceName)
				targetServiceNamePtr = &serviceName
				break
			}

			if targetServiceNamePtr != nil {
				break
			}
		}
	}

	mappingIdInt := locationBlockIndex + 1
	mappingId, err := valueObject.NewMappingId(mappingIdInt)
	if err != nil {
		return mapping, err
	}

	return entity.NewMapping(
		mappingId,
		vhostHost,
		path,
		matchPattern,
		targetType,
		targetServiceNamePtr,
		targetUrlPtr,
		targetResponseCodePtr,
		targetInlineHtmlContentPtr,
	), nil
}

func (repo VirtualHostQueryRepo) getVirtualHostMappings(
	vhost entity.VirtualHost,
) ([]entity.Mapping, error) {
	mappings := []entity.Mapping{}

	if vhost.Type.String() == "alias" {
		return mappings, nil
	}

	vhostName := vhost.Hostname
	mappingFilePath, err := repo.GetVirtualHostMappingsFilePath(vhostName)
	if err != nil {
		return mappings, err
	}

	fileContent, err := infraHelper.GetFileContent(mappingFilePath.String())
	if err != nil || len(fileContent) == 0 {
		return mappings, nil
	}

	locationBlocksRegex := regexp.MustCompile(
		`(?m)^\s*location\s(?P<modifierAndPath>.+)\s{(?P<content>[\s\S]*?\n)}`,
	)
	locationBlocks := locationBlocksRegex.FindAllStringSubmatch(fileContent, -1)
	if len(locationBlocks) == 0 {
		return mappings, errors.New("GetLocationsBlockFailed")
	}

	servicesList, err := servicesInfra.ServicesQueryRepo{}.Get()
	if err != nil {
		return mappings, errors.New("GetServicesListFailed")
	}

	serviceNamePortsMap := map[string][]string{}
	for _, service := range servicesList {
		svcNameStr := service.Name.String()
		svcPorts := []string{}
		for _, portBinding := range service.PortBindings {
			svcPorts = append(svcPorts, portBinding.Port.String())
		}

		serviceNamePortsMap[svcNameStr] = svcPorts
	}

	for locationBlockIndex, locationBlockContent := range locationBlocks {
		mapping, err := repo.locationBlockToMapping(
			locationBlockIndex,
			locationBlockContent,
			vhostName,
			serviceNamePortsMap,
		)
		if err != nil {
			log.Printf("[LocationIndex: %d] %s", locationBlockIndex, err.Error())
			continue
		}

		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

func (repo VirtualHostQueryRepo) GetWithMappings() ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings := []dto.VirtualHostWithMappings{}

	vhosts, err := repo.Get()
	if err != nil {
		return vhostsWithMappings, err
	}

	for _, vhost := range vhosts {
		mappings, err := repo.getVirtualHostMappings(vhost)
		if err != nil {
			log.Printf(
				"[%s] GetMappingsFailed: %s",
				vhost.Hostname.String(),
				err.Error(),
			)
		}

		vhostsWithMappings = append(
			vhostsWithMappings,
			dto.NewVirtualHostWithMappings(
				vhost,
				mappings,
			),
		)
	}

	return vhostsWithMappings, nil
}

func (repo VirtualHostQueryRepo) GetMappingsByHostname(
	hostname valueObject.Fqdn,
) ([]entity.Mapping, error) {
	vhostMappings := []entity.Mapping{}

	vhost, err := repo.GetByHostname(hostname)
	if err != nil {
		return vhostMappings, err
	}

	return repo.getVirtualHostMappings(vhost)
}

func (repo VirtualHostQueryRepo) GetMappingById(
	vhostHostname valueObject.Fqdn,
	id valueObject.MappingId,
) (entity.Mapping, error) {
	var mapping entity.Mapping

	vhost, err := repo.GetByHostname(vhostHostname)
	if err != nil {
		return mapping, err
	}

	mappings, err := repo.getVirtualHostMappings(vhost)
	if err != nil {
		return mapping, err
	}

	for _, mapping := range mappings {
		if mapping.Id == id {
			return mapping, nil
		}
	}

	return mapping, errors.New("MappingNotFound")
}

func (repo VirtualHostQueryRepo) CheckDomainOwnership(
	vhost valueObject.Fqdn,
	ownershipHash string,
) bool {
	ownershipValidateUrl := "https://" + vhost.String() +
		envDataInfra.DomainOwnershipValidationUrlPath

	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	httpResponse, err := httpClient.Get(ownershipValidateUrl)
	if err != nil {
		return false
	}
	defer httpResponse.Body.Close()

	responseBodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return false
	}

	achievedOwnershipHash := string(responseBodyBytes)
	return achievedOwnershipHash == ownershipHash
}
