package infra

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	"golang.org/x/exp/slices"
	"golang.org/x/net/publicsuffix"
)

var configurationsDir string = "/app/conf/nginx"
var mappingsDir string = "/app/conf/nginx/mapping"

type VirtualHostQueryRepo struct {
}

func (repo VirtualHostQueryRepo) IsVirtualHostPrimaryDomain(
	domain valueObject.Fqdn,
) bool {
	primaryDomain, err := infraHelper.GetPrimaryHostname()
	if err != nil {
		return false
	}

	return domain == primaryDomain
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
	isPrimaryDomain := repo.IsVirtualHostPrimaryDomain(firstDomain)

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

	if repo.IsVirtualHostPrimaryDomain(vhostName) {
		mappingFileName = "primary.conf"
	}

	return valueObject.NewUnixFilePath(mappingsDir + "/" + mappingFileName)
}

func (repo VirtualHostQueryRepo) locationBlockToMapping(
	locationBlockIndex int,
	locationBlockParts []string,
	vhostHost valueObject.Fqdn,
	servicesList []entity.Service,
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

	blockContent := locationBlockParts[2]
	blockContentFirstLine := strings.Split(blockContent, "\n")[0]
	blockContentFirstLineParts := strings.Split(blockContentFirstLine, " ")
	if len(blockContentFirstLineParts) == 0 {
		return mapping, errors.New("GetLocationBlockContentPartsFailed")
	}

	directive := blockContentFirstLineParts[0]
	if len(directive) == 0 {
		return mapping, errors.New("GetLocationDirectiveFailed")
	}

	targetTypeStr := "service"
	isReturn := directive == "return"
	var targetUrlPtr *valueObject.Url
	var targetResponseCodePtr *valueObject.HttpResponseCode
	if isReturn {
		responseCodeStr := blockContentFirstLineParts[1]
		if len(responseCodeStr) == 0 {
			return mapping, errors.New("InvalidReturnResponseCode: " + responseCodeStr)
		}

		responseCode, err := valueObject.NewHttpResponseCode(responseCodeStr)
		if err != nil {
			return mapping, errors.New("InvalidReturnResponseCode: " + responseCodeStr)
		}
		targetResponseCodePtr = &responseCode

		targetTypeStr = "response-code"

		hasUrl := len(blockContentFirstLineParts) == 3
		if hasUrl {
			targetTypeStr = "url"

			urlStr := blockContentFirstLineParts[2]
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
	isService := directive == "proxy_pass"
	if isService {
		serviceUrlStr := blockContentFirstLineParts[1]
		serviceUrl, err := valueObject.NewUrl(serviceUrlStr)
		if err != nil {
			return mapping, errors.New("InvalidServiceProxyPassUrl: " + serviceUrlStr)
		}

		servicePort, err := serviceUrl.GetPort()
		if err != nil {
			return mapping, errors.New("InvalidServicePort: " + serviceUrlStr)
		}

		for _, service := range servicesList {
			if !slices.Contains(service.Ports, servicePort) {
				continue
			}
			targetServiceNamePtr = &service.Name
		}

		if targetServiceNamePtr == nil {
			return mapping, errors.New("ServiceNotFound: " + serviceUrlStr)
		}
	}

	mappingId, err := valueObject.NewMappingId(locationBlockIndex)
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
		`(?m)^\s*location\s(?P<modifierAndPath>.+)\s{\n\s+(?P<content>[^}]+)*;\n}`,
	)
	locationBlocks := locationBlocksRegex.FindAllStringSubmatch(fileContent, -1)
	if len(locationBlocks) == 0 {
		return mappings, errors.New("GetLocationsBlockFailed")
	}

	servicesList, err := servicesInfra.ServicesQueryRepo{}.Get()
	if err != nil {
		return mappings, errors.New("GetServicesListFailed")
	}

	for locationBlockIndex, locationBlockContent := range locationBlocks {
		mapping, err := repo.locationBlockToMapping(
			locationBlockIndex,
			locationBlockContent,
			vhostName,
			servicesList,
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
