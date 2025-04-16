package service

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type VirtualHostService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	trailDbSvc            *internalDbInfra.TrailDatabaseService
	vhostQueryRepo        *vhostInfra.VirtualHostQueryRepo
	vhostCmdRepo          *vhostInfra.VirtualHostCmdRepo
	mappingQueryRepo      *vhostInfra.MappingQueryRepo
	mappingCmdRepo        *vhostInfra.MappingCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewVirtualHostService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostService {
	return &VirtualHostService{
		persistentDbSvc:       persistentDbSvc,
		trailDbSvc:            trailDbSvc,
		vhostQueryRepo:        vhostInfra.NewVirtualHostQueryRepo(persistentDbSvc),
		vhostCmdRepo:          vhostInfra.NewVirtualHostCmdRepo(persistentDbSvc),
		mappingQueryRepo:      vhostInfra.NewMappingQueryRepo(persistentDbSvc),
		mappingCmdRepo:        vhostInfra.NewMappingCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *VirtualHostService) VirtualHostReadRequestFactory(
	serviceInput map[string]interface{},
	withMappings bool,
) (readRequestDto dto.ReadVirtualHostsRequest, err error) {
	var hostnamePtr *valueObject.Fqdn
	if serviceInput["hostname"] != nil {
		hostname, err := valueObject.NewFqdn(serviceInput["hostname"])
		if err != nil {
			return readRequestDto, err
		}
		hostnamePtr = &hostname
	}

	var typePtr *valueObject.VirtualHostType
	if serviceInput["type"] != nil {
		vhostType, err := valueObject.NewVirtualHostType(serviceInput["type"])
		if err != nil {
			return readRequestDto, err
		}
		typePtr = &vhostType
	}

	var rootDirectoryPtr *valueObject.UnixFilePath
	if serviceInput["rootDirectory"] != nil {
		rootDirectory, err := valueObject.NewUnixFilePath(serviceInput["rootDirectory"])
		if err != nil {
			return readRequestDto, err
		}
		rootDirectoryPtr = &rootDirectory
	}

	var parentHostnamePtr *valueObject.Fqdn
	if serviceInput["parentHostname"] != nil {
		parentHostname, err := valueObject.NewFqdn(serviceInput["parentHostname"])
		if err != nil {
			return readRequestDto, err
		}
		parentHostnamePtr = &parentHostname
	}

	if serviceInput["withMappings"] != nil {
		withMappings, err = voHelper.InterfaceToBool(serviceInput["withMappings"])
		if err != nil {
			return readRequestDto, err
		}
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.VirtualHostsDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadVirtualHostsRequest{
		Pagination:      requestPagination,
		Hostname:        hostnamePtr,
		VirtualHostType: typePtr,
		RootDirectory:   rootDirectoryPtr,
		ParentHostname:  parentHostnamePtr,
		WithMappings:    &withMappings,
		CreatedBeforeAt: timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:  timeParamPtrs["createdAfterAt"],
	}, nil
}

func (service *VirtualHostService) Read(
	serviceInput map[string]interface{},
) ServiceOutput {
	readRequestDto, err := service.VirtualHostReadRequestFactory(serviceInput, false)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	readResponseDto, err := useCase.ReadVirtualHosts(service.vhostQueryRepo, readRequestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, readResponseDto)
}

func (service *VirtualHostService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"hostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	vhostType := valueObject.VirtualHostTypeTopLevel
	if input["type"] != nil {
		vhostType, err = valueObject.NewVirtualHostType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	isWildcard := false
	if input["isWildcard"] != nil {
		isWildcard, err = voHelper.InterfaceToBool(input["isWildcard"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	var parentHostnamePtr *valueObject.Fqdn
	if input["parentHostname"] != nil {
		parentHostname, err := valueObject.NewFqdn(input["parentHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		parentHostnamePtr = &parentHostname
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateVirtualHost(
		hostname, vhostType, &isWildcard, parentHostnamePtr,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateVirtualHost(
		service.vhostQueryRepo, service.vhostCmdRepo, service.activityRecordCmdRepo,
		createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "VirtualHostCreated")
}

func (service *VirtualHostService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"hostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var isWildcardPtr *bool
	if input["isWildcard"] != nil {
		isWildcard, err := voHelper.InterfaceToBool(input["isWildcard"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidIsWildcard"))
		}
		isWildcardPtr = &isWildcard
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateVirtualHost(
		hostname, isWildcardPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateVirtualHost(
		service.vhostQueryRepo, service.vhostCmdRepo, service.activityRecordCmdRepo,
		updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "VirtualHostUpdated")
}

func (service *VirtualHostService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"hostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteVirtualHost(hostname, operatorAccountId, operatorIpAddress)
	err = useCase.DeleteVirtualHost(
		service.vhostQueryRepo, service.vhostCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "VirtualHostDeleted")
}

func (service *VirtualHostService) ReadWithMappings(
	serviceInput map[string]interface{},
) ServiceOutput {
	readRequestDto, err := service.VirtualHostReadRequestFactory(serviceInput, true)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	readResponseDto, err := useCase.ReadVirtualHosts(service.vhostQueryRepo, readRequestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, readResponseDto)
}

func (service *VirtualHostService) CreateMapping(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"hostname", "path", "targetType"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	path, err := valueObject.NewMappingPath(input["path"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	matchPattern := valueObject.MappingMatchPatternBeginsWith
	if input["matchPattern"] != nil {
		matchPattern, err = valueObject.NewMappingMatchPattern(input["matchPattern"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	targetType, err := valueObject.NewMappingTargetType(input["targetType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if input["targetValue"] != nil {
		targetValue, err := valueObject.NewMappingTargetValue(
			input["targetValue"], targetType,
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if input["targetHttpResponseCode"] != nil {
		if input["targetHttpResponseCode"] == "" {
			input["targetHttpResponseCode"] = 301
		}
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			input["targetHttpResponseCode"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var shouldUpgradeInsecureRequestsPtr *bool
	if input["shouldUpgradeInsecureRequests"] != nil {
		shouldUpgradeInsecureRequests, err := voHelper.InterfaceToBool(
			input["shouldUpgradeInsecureRequests"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidShouldUpgradeInsecureRequests")
		}
		shouldUpgradeInsecureRequestsPtr = &shouldUpgradeInsecureRequests
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateMapping(
		hostname, path, matchPattern, targetType, targetValuePtr,
		targetHttpResponseCodePtr, shouldUpgradeInsecureRequestsPtr,
		nil, operatorAccountId, operatorIpAddress,
	)

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMapping(
		service.vhostQueryRepo, service.mappingCmdRepo, servicesQueryRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingCreated")
}

func (service *VirtualHostService) DeleteMapping(
	input map[string]interface{},
) ServiceOutput {
	if input["mappingId"] == nil && input["id"] != nil {
		input["mappingId"] = input["id"]
	}

	requiredParams := []string{"mappingId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	id, err := valueObject.NewMappingId(input["mappingId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteMapping(id, operatorAccountId, operatorIpAddress)
	err = useCase.DeleteMapping(
		service.mappingQueryRepo, service.mappingCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingDeleted")
}

func (service *VirtualHostService) MappingSecurityRuleReadRequestFactory(
	serviceInput map[string]interface{},
) (readRequestDto dto.ReadMappingSecurityRulesRequest, err error) {
	var mappingSecurityRuleIdPtr *valueObject.MappingSecurityRuleId
	if serviceInput["id"] != nil {
		id, err := valueObject.NewMappingSecurityRuleId(serviceInput["id"])
		if err != nil {
			return readRequestDto, err
		}
		mappingSecurityRuleIdPtr = &id
	}

	var mappingSecurityRuleNamePtr *valueObject.MappingSecurityRuleName
	if serviceInput["name"] != nil {
		name, err := valueObject.NewMappingSecurityRuleName(serviceInput["name"])
		if err != nil {
			return readRequestDto, err
		}
		mappingSecurityRuleNamePtr = &name
	}

	var allowedIpPtr *valueObject.IpAddress
	if serviceInput["allowedIp"] != nil {
		allowedIp, err := valueObject.NewIpAddress(serviceInput["allowedIp"])
		if err != nil {
			return readRequestDto, err
		}
		allowedIpPtr = &allowedIp
	}

	var blockedIpPtr *valueObject.IpAddress
	if serviceInput["blockedIp"] != nil {
		blockedIp, err := valueObject.NewIpAddress(serviceInput["blockedIp"])
		if err != nil {
			return readRequestDto, err
		}
		blockedIpPtr = &blockedIp
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.MappingSecurityRulesDefaultPagination,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadMappingSecurityRulesRequest{
		Pagination:              requestPagination,
		MappingSecurityRuleId:   mappingSecurityRuleIdPtr,
		MappingSecurityRuleName: mappingSecurityRuleNamePtr,
		AllowedIp:               allowedIpPtr,
		BlockedIp:               blockedIpPtr,
		CreatedBeforeAt:         timeParamPtrs["createdBeforeAt"],
		CreatedAfterAt:          timeParamPtrs["createdAfterAt"],
	}, nil
}

func (service *VirtualHostService) ReadMappingSecurityRules(
	serviceInput map[string]interface{},
) ServiceOutput {
	readRequestDto, err := service.MappingSecurityRuleReadRequestFactory(serviceInput)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	readResponseDto, err := useCase.ReadMappingSecurityRules(
		service.mappingQueryRepo, readRequestDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, readResponseDto)
}

func (service *VirtualHostService) CreateMappingSecurityRule(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"name"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	name, err := valueObject.NewMappingSecurityRuleName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var descriptionPtr *valueObject.MappingSecurityRuleDescription
	if input["description"] != nil {
		description, err := valueObject.NewMappingSecurityRuleDescription(input["description"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		descriptionPtr = &description
	}

	allowedIps := []valueObject.IpAddress{}
	if input["allowedIps"] != nil {
		allowedIpsInput, assertOk := input["allowedIps"].([]valueObject.IpAddress)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidAllowedIps")
		}
		allowedIps = allowedIpsInput
	}

	blockedIps := []valueObject.IpAddress{}
	if input["blockedIps"] != nil {
		blockedIpsInput, assertOk := input["blockedIps"].([]valueObject.IpAddress)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidBlockedIps")
		}
		blockedIps = blockedIpsInput
	}

	var rpsSoftLimitPerIpPtr *uint
	if input["rpsSoftLimitPerIp"] != nil {
		softLimit, err := tkVoUtil.InterfaceToUint(input["rpsSoftLimitPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRpsSoftLimitPerIp")
		}
		rpsSoftLimitPerIpPtr = &softLimit
	}

	var rpsHardLimitPerIpPtr *uint
	if input["rpsHardLimitPerIp"] != nil {
		hardLimit, err := tkVoUtil.InterfaceToUint(input["rpsHardLimitPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRpsHardLimitPerIp")
		}
		rpsHardLimitPerIpPtr = &hardLimit
	}

	var responseCodeOnMaxRequestsPtr *uint
	if input["responseCodeOnMaxRequests"] != nil {
		responseCode, err := tkVoUtil.InterfaceToUint(input["responseCodeOnMaxRequests"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidResponseCodeOnMaxRequests")
		}
		responseCodeOnMaxRequestsPtr = &responseCode
	}

	var maxConnectionsPerIpPtr *uint
	if input["maxConnectionsPerIp"] != nil {
		maxConns, err := tkVoUtil.InterfaceToUint(input["maxConnectionsPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConnectionsPerIp")
		}
		maxConnectionsPerIpPtr = &maxConns
	}

	var bandwidthBpsLimitPerConnectionPtr *valueObject.Byte
	if input["bandwidthBpsLimitPerConnection"] != nil {
		bandwidthBpsLimit, err := valueObject.NewByte(input["bandwidthBpsLimitPerConnection"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidBandwidthBpsLimitPerConnection")
		}
		bandwidthBpsLimitPerConnectionPtr = &bandwidthBpsLimit
	}

	var bandwidthLimitOnlyAfterBytesPtr *valueObject.Byte
	if input["bandwidthLimitOnlyAfterBytes"] != nil {
		bandwidthLimitOnlyAfterBytes, err := valueObject.NewByte(input["bandwidthLimitOnlyAfterBytes"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidBandwidthLimitOnlyAfterBytes")
		}
		bandwidthLimitOnlyAfterBytesPtr = &bandwidthLimitOnlyAfterBytes
	}

	var responseCodeOnMaxConnectionsPtr *uint
	if input["responseCodeOnMaxConnections"] != nil {
		responseCode, err := tkVoUtil.InterfaceToUint(input["responseCodeOnMaxConnections"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidResponseCodeOnMaxConnections")
		}
		responseCodeOnMaxConnectionsPtr = &responseCode
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateMappingSecurityRule(
		name, descriptionPtr, allowedIps, blockedIps, rpsSoftLimitPerIpPtr,
		rpsHardLimitPerIpPtr, responseCodeOnMaxRequestsPtr, maxConnectionsPerIpPtr,
		bandwidthBpsLimitPerConnectionPtr, bandwidthLimitOnlyAfterBytesPtr,
		responseCodeOnMaxConnectionsPtr, operatorAccountId, operatorIpAddress,
	)

	mappingSecurityRuleId, err := useCase.CreateMappingSecurityRule(
		service.mappingQueryRepo, service.mappingCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, map[string]interface{}{
		"id": mappingSecurityRuleId.Uint64(),
	})
}

func (service *VirtualHostService) UpdateMappingSecurityRule(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	id, err := valueObject.NewMappingSecurityRuleId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var namePtr *valueObject.MappingSecurityRuleName
	if input["name"] != nil {
		name, err := valueObject.NewMappingSecurityRuleName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		namePtr = &name
	}

	var descriptionPtr *valueObject.MappingSecurityRuleDescription
	if input["description"] != nil {
		description, err := valueObject.NewMappingSecurityRuleDescription(input["description"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		descriptionPtr = &description
	}

	allowedIps := []valueObject.IpAddress{}
	if input["allowedIps"] != nil {
		var assertOk bool
		allowedIps, assertOk = input["allowedIps"].([]valueObject.IpAddress)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidAllowedIps")
		}
	}

	blockedIps := []valueObject.IpAddress{}
	if input["blockedIps"] != nil {
		var assertOk bool
		blockedIps, assertOk = input["blockedIps"].([]valueObject.IpAddress)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidBlockedIps")
		}
	}

	var rpsSoftLimitPerIpPtr *uint
	if input["rpsSoftLimitPerIp"] != nil {
		softLimit, err := tkVoUtil.InterfaceToUint(input["rpsSoftLimitPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRpsSoftLimitPerIp")
		}
		rpsSoftLimitPerIpPtr = &softLimit
	}

	var rpsHardLimitPerIpPtr *uint
	if input["rpsHardLimitPerIp"] != nil {
		hardLimit, err := tkVoUtil.InterfaceToUint(input["rpsHardLimitPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidRpsHardLimitPerIp")
		}
		rpsHardLimitPerIpPtr = &hardLimit
	}

	var responseCodeOnMaxRequestsPtr *uint
	if input["responseCodeOnMaxRequests"] != nil {
		responseCode, err := tkVoUtil.InterfaceToUint(input["responseCodeOnMaxRequests"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidResponseCodeOnMaxRequests")
		}
		responseCodeOnMaxRequestsPtr = &responseCode
	}

	var maxConnectionsPerIpPtr *uint
	if input["maxConnectionsPerIp"] != nil {
		maxConns, err := tkVoUtil.InterfaceToUint(input["maxConnectionsPerIp"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMaxConnectionsPerIp")
		}
		maxConnectionsPerIpPtr = &maxConns
	}

	var bandwidthBpsLimitPerConnectionPtr *valueObject.Byte
	if input["bandwidthBpsLimitPerConnection"] != nil {
		bandwidthBpsLimit, err := valueObject.NewByte(input["bandwidthBpsLimitPerConnection"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidBandwidthBpsLimitPerConnection")
		}
		bandwidthBpsLimitPerConnectionPtr = &bandwidthBpsLimit
	}

	var bandwidthLimitOnlyAfterBytesPtr *valueObject.Byte
	if input["bandwidthLimitOnlyAfterBytes"] != nil {
		bandwidthLimitOnlyAfterBytes, err := valueObject.NewByte(input["bandwidthLimitOnlyAfterBytes"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidBandwidthLimitOnlyAfterBytes")
		}
		bandwidthLimitOnlyAfterBytesPtr = &bandwidthLimitOnlyAfterBytes
	}

	var responseCodeOnMaxConnectionsPtr *uint
	if input["responseCodeOnMaxConnections"] != nil {
		responseCode, err := tkVoUtil.InterfaceToUint(input["responseCodeOnMaxConnections"])
		if err != nil {
			return NewServiceOutput(UserError, "InvalidResponseCodeOnMaxConnections")
		}
		responseCodeOnMaxConnectionsPtr = &responseCode
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateMappingSecurityRule(
		id, namePtr, descriptionPtr, allowedIps, blockedIps,
		rpsSoftLimitPerIpPtr, rpsHardLimitPerIpPtr, responseCodeOnMaxRequestsPtr,
		maxConnectionsPerIpPtr, bandwidthBpsLimitPerConnectionPtr,
		bandwidthLimitOnlyAfterBytesPtr, responseCodeOnMaxConnectionsPtr,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateMappingSecurityRule(
		service.mappingQueryRepo, service.mappingCmdRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MappingSecurityRuleUpdated")
}

func (service *VirtualHostService) DeleteMappingSecurityRule(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	ruleId, err := valueObject.NewMappingSecurityRuleId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteMappingSecurityRule(
		ruleId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteMappingSecurityRule(
		service.mappingQueryRepo, service.mappingCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MappingSecurityRuleDeleted")
}
