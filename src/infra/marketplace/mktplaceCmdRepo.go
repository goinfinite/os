package mktplaceInfra

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	filesInfra "github.com/speedianet/os/src/infra/files"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type MktplaceCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	queryRepo       *MktplaceQueryRepo
}

func NewMktplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MktplaceCmdRepo {
	mktplaceQueryRepo := NewMktplaceQueryRepo(persistentDbSvc)

	return &MktplaceCmdRepo{
		persistentDbSvc: persistentDbSvc,
		queryRepo:       mktplaceQueryRepo,
	}
}

func (repo *MktplaceCmdRepo) getDataFieldsAsMap(
	dataFields []valueObject.DataField,
) map[string]string {
	dataFieldMap := map[string]string{}

	for _, dataField := range dataFields {
		dataFieldMap[dataField.Key.String()] = dataField.Value.String()
	}

	return dataFieldMap
}

func (repo *MktplaceCmdRepo) moveMktplaceItemDir(
	rootDirectory valueObject.UnixFilePath,
	mktplaceItemName valueObject.MktplaceItemName,
) error {
	mktplaceItemSrcPath, _ := valueObject.NewUnixFilePath("/speedia/" + mktplaceItemName.String())

	mktplaceItemDestinationPathStr := rootDirectory.String() + mktplaceItemName.String()
	mktplaceItemDestinationPath, _ := valueObject.NewUnixFilePath(mktplaceItemDestinationPathStr)

	err := infraHelper.MakeDir(mktplaceItemDestinationPathStr)
	if err != nil && !os.IsExist(err) {
		return err
	}

	filesCmdRepo := filesInfra.FilesCmdRepo{}
	shouldOverwrite := true
	return filesCmdRepo.Move(mktplaceItemSrcPath, mktplaceItemDestinationPath, shouldOverwrite)
}

func (repo *MktplaceCmdRepo) InstallItem(
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	mktplaceCatalogItem, err := repo.queryRepo.GetItemById(
		installMktplaceCatalogItem.Id,
	)
	if err != nil {
		return errors.New("MktplaceCatalogItemNotFound")
	}

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, requiredSvcName := range mktplaceCatalogItem.Services {
		_, err := servicesQueryRepo.GetByName(requiredSvcName)
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

		err = servicesCmdRepo.CreateInstallable(requiredService)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	dataFieldsMap := repo.getDataFieldsAsMap(installMktplaceCatalogItem.DataFields)
	for _, cmdStep := range mktplaceCatalogItem.CmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepRequiredDataField, err := infraHelper.GetRegexFirstGroup(
			cmdStepStr,
			`%(.+)%`,
		)
		if err == nil {
			requiredDataFieldValue := dataFieldsMap[cmdStepRequiredDataField]
			cmdStepWithDataField := strings.ReplaceAll(
				cmdStepStr,
				"%"+cmdStepRequiredDataField+"%",
				requiredDataFieldValue,
			)
			cmdStepStr = cmdStepWithDataField
		}

		_, err = infraHelper.RunCmdWithSubShell(cmdStepStr)
		if err != nil {
			return errors.New("RunCmdStepError (" + cmdStepStr + "): " + err.Error())
		}
	}

	err = repo.moveMktplaceItemDir(
		installMktplaceCatalogItem.RootDirectory,
		mktplaceCatalogItem.Name,
	)
	if err != nil {
		return err
	}

	nowUnixTime := time.Now().Unix()
	createdAt := valueObject.UnixTime(nowUnixTime)
	updatedAt := valueObject.UnixTime(nowUnixTime)

	mktplaceInstalledItem := entity.NewMarketplaceInstalledItem(
		mktplaceCatalogItem.Id,
		mktplaceCatalogItem.Name,
		mktplaceCatalogItem.Type,
		installMktplaceCatalogItem.RootDirectory,
		mktplaceCatalogItem.Services,
		[]entity.Mapping{},
		mktplaceCatalogItem.AvatarUrl,
		createdAt,
		updatedAt,
	)

	modelWithoutId := true
	mktplaceInstalledItemModel, err := dbModel.MarketplaceInstalledItem{}.ToModel(
		mktplaceInstalledItem,
		modelWithoutId,
	)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&mktplaceInstalledItemModel).Error
	if err != nil {
		return err
	}

	for _, mktplaceItemMapping := range mktplaceCatalogItem.Mappings {
		createMktplaceItemMapping := dto.NewCreateMapping(
			installMktplaceCatalogItem.Hostname,
			mktplaceItemMapping.Path,
			mktplaceItemMapping.MatchPattern,
			mktplaceItemMapping.TargetType,
			mktplaceItemMapping.TargetServiceName,
			mktplaceItemMapping.TargetUrl,
			mktplaceItemMapping.TargetHttpResponseCode,
			mktplaceItemMapping.TargetInlineHtmlContent,
		)

		vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
		err = vhostCmdRepo.CreateMapping(createMktplaceItemMapping)
		if err != nil {
			log.Printf("CreateMktplaceItemMappingError: %s", err.Error())
		}
	}

	return nil
}
