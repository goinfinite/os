package liaison

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
)

type RuntimeLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	availabilityInspector *sharedHelper.ServiceAvailabilityInspector
	runtimeQueryRepo      runtimeInfra.RuntimeQueryRepo
	runtimeCmdRepo        *runtimeInfra.RuntimeCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewRuntimeLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimeLiaison {
	return &RuntimeLiaison{
		persistentDbSvc: persistentDbSvc,
		availabilityInspector: sharedHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
		runtimeQueryRepo:      runtimeInfra.RuntimeQueryRepo{},
		runtimeCmdRepo:        runtimeInfra.NewRuntimeCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *RuntimeLiaison) ReadPhpConfigs(
	untrustedInput map[string]any,
) LiaisonOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	hostname, err := valueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	phpConfigs, err := useCase.ReadPhpConfigs(liaison.runtimeQueryRepo, hostname)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, phpConfigs)
}

func (liaison *RuntimeLiaison) UpdatePhpConfigs(
	untrustedInput map[string]any,
) LiaisonOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requiredParams := []string{"hostname", "version"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	phpVersion, err := valueObject.NewPhpVersion(untrustedInput["version"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	phpModules := []entity.PhpModule{}
	if _, exists := untrustedInput["modules"]; exists {
		var assertOk bool
		phpModules, assertOk = untrustedInput["modules"].([]entity.PhpModule)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidPhpModules")
		}
	}

	phpSettings := []entity.PhpSetting{}
	if _, exists := untrustedInput["settings"]; exists {
		var assertOk bool
		phpSettings, assertOk = untrustedInput["settings"].([]entity.PhpSetting)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidPhpSettings")
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdatePhpConfigs(
		hostname, phpVersion, phpModules, phpSettings, operatorAccountId,
		operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)

	err = useCase.UpdatePhpConfigs(
		liaison.runtimeQueryRepo, liaison.runtimeCmdRepo, vhostQueryRepo,
		liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "PhpConfigsUpdated")
}
