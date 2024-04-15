package marketplaceInfra

import (
	"errors"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"github.com/speedianet/os/src/infra/infraData"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type MarketplaceCmdRepo struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	marketplaceQueryRepo *MarketplaceQueryRepo
}

func NewMarketplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceCmdRepo {
	marketplaceQueryRepo := NewMarketplaceQueryRepo(persistentDbSvc)

	return &MarketplaceCmdRepo{
		persistentDbSvc:      persistentDbSvc,
		marketplaceQueryRepo: marketplaceQueryRepo,
	}
}

func (repo *MarketplaceCmdRepo) createRequiredServices(
	catalogRequiredSvcNames []valueObject.ServiceName,
) error {
	svcQueryRepo := servicesInfra.ServicesQueryRepo{}
	svcCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, requiredSvcName := range catalogRequiredSvcNames {
		_, err := svcQueryRepo.GetByName(requiredSvcName)
		if err == nil {
			continue
		}

		requiredSvcAutoCreateMapping := false
		requiredService := dto.NewCreateInstallableService(
			requiredSvcName,
			nil,
			nil,
			nil,
			requiredSvcAutoCreateMapping,
		)

		err = svcCmdRepo.CreateInstallable(requiredService)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) getMissingOptionalDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
) ([]valueObject.MarketplaceInstallableItemDataField, error) {
	missingCatalogOptionalDataFields := []valueObject.MarketplaceInstallableItemDataField{}

	receivedDataFieldsKeys := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsKeys = append(
			receivedDataFieldsKeys,
			receivedDataField.Key.String(),
		)
	}

	for _, catalogDataField := range catalogDataFields {
		if catalogDataField.IsRequired {
			continue
		}

		catalogDataFieldKeyStr := catalogDataField.Key.String()
		if slices.Contains(receivedDataFieldsKeys, catalogDataFieldKeyStr) {
			continue
		}

		if catalogDataField.DefaultValue == nil {
			return missingCatalogOptionalDataFields, errors.New(
				"CatalogOptionalDataFieldWithoutDefaultValue: " +
					catalogDataFieldKeyStr,
			)
		}

		catalogDataFieldAsInstallable, _ := valueObject.NewMarketplaceInstallableItemDataField(
			catalogDataField.Key,
			*catalogDataField.DefaultValue,
		)
		missingCatalogOptionalDataFields = append(
			missingCatalogOptionalDataFields,
			catalogDataFieldAsInstallable,
		)
	}

	return missingCatalogOptionalDataFields, nil
}

func (repo *MarketplaceCmdRepo) replaceCmdStepsPlaceholders(
	cmdSteps []valueObject.MarketplaceItemCmdStep,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) ([]valueObject.MarketplaceItemCmdStep, error) {
	cmdStepsWithDataFields := []valueObject.MarketplaceItemCmdStep{}

	dataFieldsMap := map[string]string{}
	for _, dataField := range dataFields {
		dataFieldKeyStr := dataField.Key.String()
		dataFieldsMap[dataFieldKeyStr] = dataField.Value.String()
	}

	for _, cmdStep := range cmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepDataFieldKeys, _ := infraHelper.GetAllRegexGroupMatches(
			cmdStepStr,
			`%(.*?)%`,
		)

		for _, cmdStepDataFieldKey := range cmdStepDataFieldKeys {
			dataFieldValue := dataFieldsMap[cmdStepDataFieldKey]
			cmdStepWithDataFieldStr := strings.ReplaceAll(
				cmdStepStr,
				"%"+cmdStepDataFieldKey+"%",
				dataFieldValue,
			)
			cmdStepStr = cmdStepWithDataFieldStr
		}

		cmdStepWithDataField, _ := valueObject.NewMarketplaceItemCmdStep(cmdStepStr)
		cmdStepsWithDataFields = append(cmdStepsWithDataFields, cmdStepWithDataField)
	}

	return cmdStepsWithDataFields, nil
}

func (repo *MarketplaceCmdRepo) runCmdSteps(
	catalogCmdSteps []valueObject.MarketplaceItemCmdStep,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	installDir valueObject.UnixFilePath,
	installUuid string,
) error {
	installDirDataFieldKey, _ := valueObject.NewDataFieldKey("installDirectory")
	installDirDataFieldValue, _ := valueObject.NewDataFieldValue(installDir.String())
	installDirDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installDirDataFieldKey,
		installDirDataFieldValue,
	)

	installUuidDataFieldKey, _ := valueObject.NewDataFieldKey("installUuid")
	installUuidDataFieldValue, _ := valueObject.NewDataFieldValue(installUuid)
	installUuidDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installUuidDataFieldKey,
		installUuidDataFieldValue,
	)

	receivedDataFields = append(
		receivedDataFields,
		installDirDataField,
		installUuidDataField,
	)

	missingCatalogOptionalDataFields, err := repo.getMissingOptionalDataFields(
		receivedDataFields,
		catalogDataFields,
	)
	if err != nil {
		return err
	}
	receivedDataFields = slices.Concat(
		receivedDataFields,
		missingCatalogOptionalDataFields,
	)

	log.Printf("ReceivedDataFields: %+v", receivedDataFields)

	preparedCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogCmdSteps,
		receivedDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	for preparedCmdStepIndex, preparedCmdStep := range preparedCmdSteps {
		preparedCmdStepStr := preparedCmdStep.String()
		_, err = infraHelper.RunCmdWithSubShell(preparedCmdStepStr)
		if err != nil {
			preparedCmdStepIndexStr := strconv.Itoa(preparedCmdStepIndex)
			return errors.New(
				"RunCmdStepError (" + preparedCmdStepIndexStr + "): " + err.Error(),
			)
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) createMappings(
	hostname valueObject.Fqdn,
	catalogMappings []valueObject.MarketplaceItemMapping,
) error {
	for _, catalogMapping := range catalogMappings {
		createCatalogItemMapping := dto.NewCreateMapping(
			hostname,
			catalogMapping.Path,
			catalogMapping.MatchPattern,
			catalogMapping.TargetType,
			catalogMapping.TargetServiceName,
			catalogMapping.TargetUrl,
			catalogMapping.TargetHttpResponseCode,
			catalogMapping.TargetInlineHtmlContent,
		)

		vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
		err := vhostCmdRepo.CreateMapping(createCatalogItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	installDir valueObject.UnixFilePath,
	installUuid string,
) error {
	svcNamesListStr := []string{}
	for _, svcName := range catalogItem.ServiceNames {
		svcNamesListStr = append(svcNamesListStr, svcName.String())
	}
	svcNamesStr := strings.Join(svcNamesListStr, ",")

	installedItemModel := dbModel.MarketplaceInstalledItem{
		Name:             catalogItem.Name.String(),
		Type:             catalogItem.Type.String(),
		InstallDirectory: installDir.String(),
		InstallUuid:      installUuid,
		ServiceNames:     svcNamesStr,
		AvatarUrl:        catalogItem.AvatarUrl.String(),
	}

	err := repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	catalogItem, err := repo.marketplaceQueryRepo.GetCatalogItemById(
		installDto.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	err = repo.createRequiredServices(catalogItem.ServiceNames)
	if err != nil {
		return err
	}

	installDirStr := infraData.GlobalConfigs.PrimaryPublicDir
	if installDto.InstallDirectory != nil {
		vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
		vhost, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
		if err != nil {
			return err
		}

		installDirStr = installDto.InstallDirectory.String()
		hasLeadingSlash := strings.HasPrefix(installDirStr, "/")
		if !hasLeadingSlash {
			installDirStr = "/" + installDirStr
		}

		installDirStr = vhost.RootDirectory.String() + installDirStr
	}
	installDir, _ := valueObject.NewUnixFilePath(installDirStr)

	installUuid := uuid.New().String()[:16]

	err = repo.runCmdSteps(
		catalogItem.CmdSteps,
		catalogItem.DataFields,
		installDto.DataFields,
		installDir,
		installUuid,
	)
	if err != nil {
		return err
	}

	err = repo.createMappings(installDto.Hostname, catalogItem.Mappings)
	if err != nil {
		return err
	}

	return repo.persistInstalledItem(catalogItem, installDir, installUuid)
}
