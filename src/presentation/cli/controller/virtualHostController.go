package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type VirtualHostController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		persistentDbSvc: persistentDbSvc,
	}
}

func (controller *VirtualHostController) Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetVirtualHosts",
		Run: func(cmd *cobra.Command, args []string) {
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
			vhostsList, err := useCase.GetVirtualHosts(vhostQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsList)
		},
	}

	return cmd
}

func (controller *VirtualHostController) Create() *cobra.Command {
	var hostnameStr string
	var typeStr string
	var parentHostnameStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)

			vhostTypeStr := "top-level"
			if typeStr != "" {
				vhostTypeStr = typeStr
			}
			vhostType := valueObject.NewVirtualHostTypePanic(vhostTypeStr)

			var parentHostnamePtr *valueObject.Fqdn
			if parentHostnameStr != "" {
				parentHostname := valueObject.NewFqdnPanic(parentHostnameStr)
				parentHostnamePtr = &parentHostname
			}

			createVirtualHostDto := dto.NewCreateVirtualHost(
				hostname,
				vhostType,
				parentHostnamePtr,
			)

			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
			vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

			err := useCase.CreateVirtualHost(
				vhostQueryRepo,
				vhostCmdRepo,
				createVirtualHostDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "VirtualHostCreated")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(
		&typeStr, "type", "t", "", "VirtualHostType (top-level|subdomain|alias)",
	)
	cmd.Flags().StringVarP(
		&parentHostnameStr, "parent", "p", "", "ParentHostname",
	)
	return cmd
}

func (controller *VirtualHostController) Delete() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)

			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
			vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

			primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
			if err != nil {
				panic("PrimaryVirtualHostNotFound")
			}

			err = useCase.DeleteVirtualHost(
				vhostQueryRepo,
				vhostCmdRepo,
				primaryVhost,
				hostname,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "VirtualHostDeleted")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	return cmd
}

func (controller *VirtualHostController) GetWithMappings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetVirtualHostsWithMappings",
		Run: func(cmd *cobra.Command, args []string) {
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(
				controller.persistentDbSvc,
			)

			vhostsWithMappings, err := useCase.ReadVirtualHostsWithMappings(
				mappingQueryRepo,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsWithMappings)
		},
	}

	return cmd
}

func (controller *VirtualHostController) CreateMapping() *cobra.Command {
	var hostnameStr string
	var pathStr string
	var matchPatternStr string
	var targetTypeStr string
	var targetValueStr string
	var targetHttpResponseCodeUint uint

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateVirtualHostMapping",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			path := valueObject.NewMappingPathPanic(pathStr)

			matchPattern := valueObject.NewMappingMatchPatternPanic("begins-with")
			if matchPatternStr != "" {
				matchPattern = valueObject.NewMappingMatchPatternPanic(matchPatternStr)
			}

			targetType := valueObject.NewMappingTargetTypePanic(targetTypeStr)

			var targetValuePtr *valueObject.MappingTargetValue
			if targetValueStr != "" {
				targetValue := valueObject.NewMappingTargetValuePanic(
					targetValueStr, targetType,
				)
				targetValuePtr = &targetValue
			}

			var targetHttpResponseCodePtr *valueObject.HttpResponseCode
			if targetHttpResponseCodeUint != 0 {
				targetHttpResponseCode := valueObject.NewHttpResponseCodePanic(
					targetHttpResponseCodeUint,
				)
				targetHttpResponseCodePtr = &targetHttpResponseCode
			}

			createMappingDto := dto.NewCreateMapping(
				hostname,
				path,
				matchPattern,
				targetType,
				targetValuePtr,
				targetHttpResponseCodePtr,
			)

			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
			svcsQueryRepo := servicesInfra.ServicesQueryRepo{}

			err := useCase.CreateMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				vhostQueryRepo,
				svcsQueryRepo,
				createMappingDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingCreated")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&pathStr, "path", "p", "", "MappingPath")
	cmd.MarkFlagRequired("path")
	cmd.Flags().StringVarP(
		&matchPatternStr, "match", "m", "",
		"MatchPattern (begins-with|contains|ends-with)",
	)
	cmd.Flags().StringVarP(
		&targetTypeStr, "type", "t", "",
		"MappingTargetType (url|service|response-code|inline-html|static-files)",
	)
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&targetValueStr, "value", "v", "", "MappingTargetValue")
	cmd.Flags().UintVarP(
		&targetHttpResponseCodeUint, "response-code", "r", 0, "TargetHttpResponseCode",
	)
	return cmd
}

func (controller *VirtualHostController) DeleteMapping() *cobra.Command {
	var mappingIdUint uint

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHostMapping",
		Run: func(cmd *cobra.Command, args []string) {
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)

			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				mappingId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingDeleted")
		},
	}

	cmd.Flags().UintVarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}
