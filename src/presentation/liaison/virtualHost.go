package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type VirtualHostLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	trailDbSvc            *internalDbInfra.TrailDatabaseService
	vhostQueryRepo        *vhostInfra.VirtualHostQueryRepo
	vhostCmdRepo          *vhostInfra.VirtualHostCmdRepo
	mappingQueryRepo      *vhostInfra.MappingQueryRepo
	mappingCmdRepo        *vhostInfra.MappingCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewVirtualHostLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostLiaison {
	return &VirtualHostLiaison{
		persistentDbSvc:       persistentDbSvc,
		trailDbSvc:            trailDbSvc,
		vhostQueryRepo:        vhostInfra.NewVirtualHostQueryRepo(persistentDbSvc),
		vhostCmdRepo:          vhostInfra.NewVirtualHostCmdRepo(persistentDbSvc),
		mappingQueryRepo:      vhostInfra.NewMappingQueryRepo(persistentDbSvc),
		mappingCmdRepo:        vhostInfra.NewMappingCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *VirtualHostLiaison) VirtualHostReadRequestFactory(
	untrustedInput map[string]any,
	withMappings bool,
) (readRequestDto dto.ReadVirtualHostsRequest, err error) {
	var hostnamePtr *tkValueObject.Fqdn
	if untrustedInput["hostname"] != nil {
		hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
		if err != nil {
			return readRequestDto, err
		}
		hostnamePtr = &hostname
	}

	var typePtr *valueObject.VirtualHostType
	if untrustedInput["type"] != nil {
		vhostType, err := valueObject.NewVirtualHostType(untrustedInput["type"])
		if err != nil {
			return readRequestDto, err
		}
		typePtr = &vhostType
	}

	var rootDirectoryPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["rootDirectory"] != nil {
		rootDirectory, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["rootDirectory"], false)
		if err != nil {
			return readRequestDto, err
		}
		rootDirectoryPtr = &rootDirectory
	}

	var parentHostnamePtr *tkValueObject.Fqdn
	if untrustedInput["parentHostname"] != nil {
		parentHostname, err := tkValueObject.NewFqdn(untrustedInput["parentHostname"])
		if err != nil {
			return readRequestDto, err
		}
		parentHostnamePtr = &parentHostname
	}

	if untrustedInput["withMappings"] != nil {
		withMappings, err = tkVoUtil.InterfaceToBool(untrustedInput["withMappings"])
		if err != nil {
			return readRequestDto, err
		}
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := tkPresentation.TimeParamsParser(timeParamNames, untrustedInput)

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.VirtualHostsDefaultPagination, untrustedInput,
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

func (liaison *VirtualHostLiaison) Read(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	readRequestDto, err := liaison.VirtualHostReadRequestFactory(untrustedInput, false)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	readResponseDto, err := useCase.ReadVirtualHosts(liaison.vhostQueryRepo, readRequestDto)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		readResponseDto,
	)
}

func (liaison *VirtualHostLiaison) Create(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"hostname"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	vhostType := valueObject.VirtualHostTypeTopLevel
	if untrustedInput["type"] != nil {
		vhostType, err = valueObject.NewVirtualHostType(untrustedInput["type"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	isWildcard := false
	if untrustedInput["isWildcard"] != nil {
		isWildcard, err = tkVoUtil.InterfaceToBool(untrustedInput["isWildcard"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	var parentHostnamePtr *tkValueObject.Fqdn
	if untrustedInput["parentHostname"] != nil {
		parentHostname, err := tkValueObject.NewFqdn(untrustedInput["parentHostname"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		parentHostnamePtr = &parentHostname
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	createDto := dto.NewCreateVirtualHost(
		hostname, vhostType, &isWildcard, parentHostnamePtr,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateVirtualHost(
		liaison.vhostQueryRepo, liaison.vhostCmdRepo, liaison.activityRecordCmdRepo,
		createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated,
		"VirtualHostCreated",
	)
}

func (liaison *VirtualHostLiaison) Update(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"hostname"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	var isWildcardPtr *bool
	if untrustedInput["isWildcard"] != nil {
		isWildcard, err := tkVoUtil.InterfaceToBool(untrustedInput["isWildcard"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				errors.New("InvalidIsWildcard"),
			)
		}
		isWildcardPtr = &isWildcard
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	updateDto := dto.NewUpdateVirtualHost(
		hostname, isWildcardPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateVirtualHost(
		liaison.vhostQueryRepo, liaison.vhostCmdRepo, liaison.activityRecordCmdRepo,
		updateDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"VirtualHostUpdated",
	)
}

func (liaison *VirtualHostLiaison) Delete(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"hostname"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	deleteDto := dto.NewDeleteVirtualHost(hostname, operatorAccountId, operatorIpAddress)
	err = useCase.DeleteVirtualHost(
		liaison.vhostQueryRepo, liaison.vhostCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"VirtualHostDeleted",
	)
}

func (liaison *VirtualHostLiaison) ReadWithMappings(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	readRequestDto, err := liaison.VirtualHostReadRequestFactory(untrustedInput, true)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	readResponseDto, err := useCase.ReadVirtualHosts(liaison.vhostQueryRepo, readRequestDto)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		readResponseDto,
	)
}

func (liaison *VirtualHostLiaison) CreateMapping(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"hostname", "path", "targetType"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	path, err := valueObject.NewMappingPath(untrustedInput["path"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	matchPattern := valueObject.MappingMatchPatternBeginsWith
	if untrustedInput["matchPattern"] != nil {
		matchPattern, err = valueObject.NewMappingMatchPattern(untrustedInput["matchPattern"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	targetType, err := valueObject.NewMappingTargetType(untrustedInput["targetType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if untrustedInput["targetValue"] != nil {
		targetValue, err := valueObject.NewMappingTargetValue(
			untrustedInput["targetValue"], targetType,
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *tkValueObject.HttpStatusCode
	if untrustedInput["targetHttpResponseCode"] != nil {
		if untrustedInput["targetHttpResponseCode"] == "" {
			untrustedInput["targetHttpResponseCode"] = 301
		}
		targetHttpResponseCode, err := tkValueObject.NewHttpStatusCode(
			untrustedInput["targetHttpResponseCode"],
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var shouldUpgradeInsecureRequestsPtr *bool
	if untrustedInput["shouldUpgradeInsecureRequests"] != nil {
		shouldUpgradeInsecureRequests, err := tkVoUtil.InterfaceToBool(
			untrustedInput["shouldUpgradeInsecureRequests"],
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidShouldUpgradeInsecureRequests",
			)
		}
		shouldUpgradeInsecureRequestsPtr = &shouldUpgradeInsecureRequests
	}

	var mappingSecurityRuleIdPtr *valueObject.MappingSecurityRuleId
	if untrustedInput["mappingSecurityRuleId"] != nil && untrustedInput["mappingSecurityRuleId"] != "" {
		mappingSecurityRuleId, err := valueObject.NewMappingSecurityRuleId(
			untrustedInput["mappingSecurityRuleId"],
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		mappingSecurityRuleIdPtr = &mappingSecurityRuleId
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	createDto := dto.NewCreateMapping(
		hostname, path, matchPattern, targetType, targetValuePtr,
		targetHttpResponseCodePtr, shouldUpgradeInsecureRequestsPtr,
		mappingSecurityRuleIdPtr, operatorAccountId, operatorIpAddress,
	)

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(liaison.persistentDbSvc)
	runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(liaison.persistentDbSvc)

	err = useCase.CreateMapping(
		liaison.vhostQueryRepo, liaison.mappingCmdRepo, runtimeCmdRepo,
		servicesQueryRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated,
		"MappingCreated",
	)
}

func (liaison *VirtualHostLiaison) DeleteMapping(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	if untrustedInput["mappingId"] == nil && untrustedInput["id"] != nil {
		untrustedInput["mappingId"] = untrustedInput["id"]
	}

	requiredParams := []string{"mappingId"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	id, err := valueObject.NewMappingId(untrustedInput["mappingId"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	deleteDto := dto.NewDeleteMapping(id, operatorAccountId, operatorIpAddress)
	err = useCase.DeleteMapping(
		liaison.mappingQueryRepo, liaison.mappingCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"MappingDeleted",
	)
}

func (liaison *VirtualHostLiaison) UpdateMapping(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"id"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	id, err := valueObject.NewMappingId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	var pathPtr *valueObject.MappingPath
	if untrustedInput["path"] != nil {
		path, err := valueObject.NewMappingPath(untrustedInput["path"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		pathPtr = &path
	}

	var matchPatternPtr *valueObject.MappingMatchPattern
	if untrustedInput["matchPattern"] != nil {
		matchPattern, err := valueObject.NewMappingMatchPattern(untrustedInput["matchPattern"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		matchPatternPtr = &matchPattern
	}

	var targetTypePtr *valueObject.MappingTargetType
	if untrustedInput["targetType"] != nil {
		targetType, err := valueObject.NewMappingTargetType(untrustedInput["targetType"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		targetTypePtr = &targetType
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if untrustedInput["targetValue"] != nil {
		if targetTypePtr == nil {
			mappingEntity, err := liaison.mappingQueryRepo.ReadFirst(
				dto.ReadMappingsRequest{MappingId: &id},
			)
			if err != nil {
				return tkPresentation.NewLiaisonResponseNoMessage(
					tkPresentation.LiaisonResponseStatusInfraError,
					"ReadMappingEntityToRetrieveTargetTypeError",
				)
			}
			targetTypePtr = &mappingEntity.TargetType
		}

		targetValue, err := valueObject.NewMappingTargetValue(untrustedInput["targetValue"], *targetTypePtr)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *tkValueObject.HttpStatusCode
	if untrustedInput["targetHttpResponseCode"] != nil {
		targetHttpResponseCode, err := tkValueObject.NewHttpStatusCode(untrustedInput["targetHttpResponseCode"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var shouldUpgradeInsecureRequestsPtr *bool
	if untrustedInput["shouldUpgradeInsecureRequests"] != nil {
		shouldUpgradeInsecureRequests, err := tkVoUtil.InterfaceToBool(untrustedInput["shouldUpgradeInsecureRequests"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				errors.New("InvalidShouldUpgradeInsecureRequests"),
			)
		}
		shouldUpgradeInsecureRequestsPtr = &shouldUpgradeInsecureRequests
	}

	clearableFields := []string{}
	var mappingSecurityRuleIdPtr *valueObject.MappingSecurityRuleId
	switch mappingSecurityRuleIdValue := untrustedInput["mappingSecurityRuleId"]; {
	case mappingSecurityRuleIdValue == nil:
	case mappingSecurityRuleIdValue == "" || mappingSecurityRuleIdValue == " ":
		clearableFields = append(clearableFields, "mappingSecurityRuleId")
	default:
		mappingSecurityRuleId, err := valueObject.NewMappingSecurityRuleId(mappingSecurityRuleIdValue)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		mappingSecurityRuleIdPtr = &mappingSecurityRuleId
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	updateDto := dto.NewUpdateMapping(
		id, pathPtr, matchPatternPtr, targetTypePtr, targetValuePtr,
		targetHttpResponseCodePtr, shouldUpgradeInsecureRequestsPtr,
		mappingSecurityRuleIdPtr, clearableFields,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateMapping(
		liaison.mappingQueryRepo, liaison.mappingCmdRepo,
		liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"MappingUpdated",
	)
}

func (liaison *VirtualHostLiaison) MappingSecurityRuleReadRequestFactory(
	untrustedInput map[string]any,
) (readRequestDto dto.ReadMappingSecurityRulesRequest, err error) {
	var mappingSecurityRuleIdPtr *valueObject.MappingSecurityRuleId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewMappingSecurityRuleId(untrustedInput["id"])
		if err != nil {
			return readRequestDto, err
		}
		mappingSecurityRuleIdPtr = &id
	}

	var mappingSecurityRuleNamePtr *valueObject.MappingSecurityRuleName
	if untrustedInput["name"] != nil {
		name, err := valueObject.NewMappingSecurityRuleName(untrustedInput["name"])
		if err != nil {
			return readRequestDto, err
		}
		mappingSecurityRuleNamePtr = &name
	}

	var allowedIpPtr *tkValueObject.CidrBlock
	if untrustedInput["allowedIp"] != nil {
		allowedIp, err := tkValueObject.NewCidrBlock(untrustedInput["allowedIp"])
		if err != nil {
			return readRequestDto, err
		}
		allowedIpPtr = &allowedIp
	}

	var blockedIpPtr *tkValueObject.CidrBlock
	if untrustedInput["blockedIp"] != nil {
		blockedIp, err := tkValueObject.NewCidrBlock(untrustedInput["blockedIp"])
		if err != nil {
			return readRequestDto, err
		}
		blockedIpPtr = &blockedIp
	}

	timeParamNames := []string{"createdBeforeAt", "createdAfterAt"}
	timeParamPtrs := tkPresentation.TimeParamsParser(timeParamNames, untrustedInput)

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.MappingSecurityRulesDefaultPagination, untrustedInput,
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

func (liaison *VirtualHostLiaison) ReadMappingSecurityRules(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	readRequestDto, err := liaison.MappingSecurityRuleReadRequestFactory(untrustedInput)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	readResponseDto, err := useCase.ReadMappingSecurityRules(
		liaison.mappingQueryRepo, readRequestDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		readResponseDto,
	)
}

func (liaison *VirtualHostLiaison) CreateMappingSecurityRule(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"name"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	name, err := valueObject.NewMappingSecurityRuleName(untrustedInput["name"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	var descriptionPtr *valueObject.MappingSecurityRuleDescription
	if untrustedInput["description"] != nil {
		description, err := valueObject.NewMappingSecurityRuleDescription(untrustedInput["description"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		descriptionPtr = &description
	}

	allowedIps := []tkValueObject.CidrBlock{}
	if untrustedInput["allowedIps"] != nil {
		allowedIpsInput, assertOk := untrustedInput["allowedIps"].([]tkValueObject.CidrBlock)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidAllowedIps",
			)
		}
		allowedIps = allowedIpsInput
	}

	blockedIps := []tkValueObject.CidrBlock{}
	if untrustedInput["blockedIps"] != nil {
		blockedIpsInput, assertOk := untrustedInput["blockedIps"].([]tkValueObject.CidrBlock)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBlockedIps",
			)
		}
		blockedIps = blockedIpsInput
	}

	var rpsSoftLimitPerIpPtr *uint
	if untrustedInput["rpsSoftLimitPerIp"] != nil && untrustedInput["rpsSoftLimitPerIp"] != "" {
		softLimit, err := tkVoUtil.InterfaceToUint(untrustedInput["rpsSoftLimitPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidRpsSoftLimitPerIp",
			)
		}
		rpsSoftLimitPerIpPtr = &softLimit
	}

	var rpsHardLimitPerIpPtr *uint
	if untrustedInput["rpsHardLimitPerIp"] != nil && untrustedInput["rpsHardLimitPerIp"] != "" {
		hardLimit, err := tkVoUtil.InterfaceToUint(untrustedInput["rpsHardLimitPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidRpsHardLimitPerIp",
			)
		}
		rpsHardLimitPerIpPtr = &hardLimit
	}

	var responseCodeOnMaxRequestsPtr *uint
	if untrustedInput["responseCodeOnMaxRequests"] != nil && untrustedInput["responseCodeOnMaxRequests"] != "" {
		responseCode, err := tkVoUtil.InterfaceToUint(untrustedInput["responseCodeOnMaxRequests"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidResponseCodeOnMaxRequests",
			)
		}
		responseCodeOnMaxRequestsPtr = &responseCode
	}

	var maxConnectionsPerIpPtr *uint
	if untrustedInput["maxConnectionsPerIp"] != nil && untrustedInput["maxConnectionsPerIp"] != "" {
		maxConns, err := tkVoUtil.InterfaceToUint(untrustedInput["maxConnectionsPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidMaxConnectionsPerIp",
			)
		}
		maxConnectionsPerIpPtr = &maxConns
	}

	var bandwidthBpsLimitPerConnectionPtr *tkValueObject.Byte
	if untrustedInput["bandwidthBpsLimitPerConnection"] != nil && untrustedInput["bandwidthBpsLimitPerConnection"] != "" {
		bandwidthBpsLimit, err := tkValueObject.NewByte(untrustedInput["bandwidthBpsLimitPerConnection"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBandwidthBpsLimitPerConnection",
			)
		}
		bandwidthBpsLimitPerConnectionPtr = &bandwidthBpsLimit
	}

	var bandwidthLimitOnlyAfterBytesPtr *tkValueObject.Byte
	if untrustedInput["bandwidthLimitOnlyAfterBytes"] != nil && untrustedInput["bandwidthLimitOnlyAfterBytes"] != "" {
		bandwidthLimitOnlyAfterBytes, err := tkValueObject.NewByte(untrustedInput["bandwidthLimitOnlyAfterBytes"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBandwidthLimitOnlyAfterBytes",
			)
		}
		bandwidthLimitOnlyAfterBytesPtr = &bandwidthLimitOnlyAfterBytes
	}

	var responseCodeOnMaxConnectionsPtr *uint
	if untrustedInput["responseCodeOnMaxConnections"] != nil && untrustedInput["responseCodeOnMaxConnections"] != "" {
		responseCode, err := tkVoUtil.InterfaceToUint(untrustedInput["responseCodeOnMaxConnections"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidResponseCodeOnMaxConnections",
			)
		}
		responseCodeOnMaxConnectionsPtr = &responseCode
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	createDto := dto.NewCreateMappingSecurityRule(
		name, descriptionPtr, allowedIps, blockedIps, rpsSoftLimitPerIpPtr,
		rpsHardLimitPerIpPtr, responseCodeOnMaxRequestsPtr, maxConnectionsPerIpPtr,
		bandwidthBpsLimitPerConnectionPtr, bandwidthLimitOnlyAfterBytesPtr,
		responseCodeOnMaxConnectionsPtr, operatorAccountId, operatorIpAddress,
	)

	mappingSecurityRuleId, err := useCase.CreateMappingSecurityRule(
		liaison.mappingQueryRepo, liaison.mappingCmdRepo,
		liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated,
		map[string]interface{}{
			"id": mappingSecurityRuleId.Uint64(),
		},
	)
}

func (liaison *VirtualHostLiaison) UpdateMappingSecurityRule(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"id"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	id, err := valueObject.NewMappingSecurityRuleId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	var namePtr *valueObject.MappingSecurityRuleName
	if untrustedInput["name"] != nil {
		name, err := valueObject.NewMappingSecurityRuleName(untrustedInput["name"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		namePtr = &name
	}

	clearableFields := []string{}

	var descriptionPtr *valueObject.MappingSecurityRuleDescription
	switch descriptionValue := untrustedInput["description"]; {
	case descriptionValue == nil:
	case descriptionValue == "" || descriptionValue == " ":
		clearableFields = append(clearableFields, "description")
	default:
		description, err := valueObject.NewMappingSecurityRuleDescription(descriptionValue)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
		descriptionPtr = &description
	}

	allowedIps := []tkValueObject.CidrBlock{}
	if untrustedInput["allowedIps"] != nil {
		var assertOk bool
		allowedIps, assertOk = untrustedInput["allowedIps"].([]tkValueObject.CidrBlock)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidAllowedIps",
			)
		}
		if len(allowedIps) == 0 {
			clearableFields = append(clearableFields, "allowedIps")
		}
	}

	blockedIps := []tkValueObject.CidrBlock{}
	if untrustedInput["blockedIps"] != nil {
		var assertOk bool
		blockedIps, assertOk = untrustedInput["blockedIps"].([]tkValueObject.CidrBlock)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBlockedIps",
			)
		}
		if len(blockedIps) == 0 {
			clearableFields = append(clearableFields, "blockedIps")
		}
	}

	var rpsSoftLimitPerIpPtr *uint
	if untrustedInput["rpsSoftLimitPerIp"] != nil && untrustedInput["rpsSoftLimitPerIp"] != "" {
		softLimit, err := tkVoUtil.InterfaceToUint(untrustedInput["rpsSoftLimitPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidRpsSoftLimitPerIp",
			)
		}
		rpsSoftLimitPerIpPtr = &softLimit
	}

	var rpsHardLimitPerIpPtr *uint
	if untrustedInput["rpsHardLimitPerIp"] != nil && untrustedInput["rpsHardLimitPerIp"] != "" {
		hardLimit, err := tkVoUtil.InterfaceToUint(untrustedInput["rpsHardLimitPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidRpsHardLimitPerIp",
			)
		}
		rpsHardLimitPerIpPtr = &hardLimit
	}

	var responseCodeOnMaxRequestsPtr *uint
	if untrustedInput["responseCodeOnMaxRequests"] != nil && untrustedInput["responseCodeOnMaxRequests"] != "" {
		responseCode, err := tkVoUtil.InterfaceToUint(untrustedInput["responseCodeOnMaxRequests"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidResponseCodeOnMaxRequests",
			)
		}
		responseCodeOnMaxRequestsPtr = &responseCode
	}

	var maxConnectionsPerIpPtr *uint
	if untrustedInput["maxConnectionsPerIp"] != nil && untrustedInput["maxConnectionsPerIp"] != "" {
		maxConns, err := tkVoUtil.InterfaceToUint(untrustedInput["maxConnectionsPerIp"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidMaxConnectionsPerIp",
			)
		}
		maxConnectionsPerIpPtr = &maxConns
	}

	var bandwidthBpsLimitPerConnectionPtr *tkValueObject.Byte
	if untrustedInput["bandwidthBpsLimitPerConnection"] != nil && untrustedInput["bandwidthBpsLimitPerConnection"] != "" {
		bandwidthBpsLimit, err := tkValueObject.NewByte(untrustedInput["bandwidthBpsLimitPerConnection"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBandwidthBpsLimitPerConnection",
			)
		}
		bandwidthBpsLimitPerConnectionPtr = &bandwidthBpsLimit
	}

	var bandwidthLimitOnlyAfterBytesPtr *tkValueObject.Byte
	if untrustedInput["bandwidthLimitOnlyAfterBytes"] != nil && untrustedInput["bandwidthLimitOnlyAfterBytes"] != "" {
		bandwidthLimitOnlyAfterBytes, err := tkValueObject.NewByte(untrustedInput["bandwidthLimitOnlyAfterBytes"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidBandwidthLimitOnlyAfterBytes",
			)
		}
		bandwidthLimitOnlyAfterBytesPtr = &bandwidthLimitOnlyAfterBytes
	}

	var responseCodeOnMaxConnectionsPtr *uint
	if untrustedInput["responseCodeOnMaxConnections"] != nil && untrustedInput["responseCodeOnMaxConnections"] != "" {
		responseCode, err := tkVoUtil.InterfaceToUint(untrustedInput["responseCodeOnMaxConnections"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				"InvalidResponseCodeOnMaxConnections",
			)
		}
		responseCodeOnMaxConnectionsPtr = &responseCode
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	updateDto := dto.NewUpdateMappingSecurityRule(
		id, namePtr, descriptionPtr, allowedIps, blockedIps,
		rpsSoftLimitPerIpPtr, rpsHardLimitPerIpPtr, responseCodeOnMaxRequestsPtr,
		maxConnectionsPerIpPtr, bandwidthBpsLimitPerConnectionPtr,
		bandwidthLimitOnlyAfterBytesPtr, responseCodeOnMaxConnectionsPtr,
		clearableFields, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateMappingSecurityRule(
		liaison.mappingQueryRepo, liaison.mappingCmdRepo,
		liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"MappingSecurityRuleUpdated",
	)
}

func (liaison *VirtualHostLiaison) DeleteMappingSecurityRule(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"id"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	ruleId, err := valueObject.NewMappingSecurityRuleId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError,
			err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError,
				err.Error(),
			)
		}
	}

	deleteDto := dto.NewDeleteMappingSecurityRule(
		ruleId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteMappingSecurityRule(
		liaison.mappingQueryRepo, liaison.mappingCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError,
			err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess,
		"MappingSecurityRuleDeleted",
	)
}
