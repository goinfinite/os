package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type VirtualHostController struct {
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	virtualHostService *service.VirtualHostService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		persistentDbSvc:    persistentDbSvc,
		virtualHostService: service.NewVirtualHostService(persistentDbSvc),
	}
}

func (controller *VirtualHostController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadVirtualHosts",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.virtualHostService.Read())
		},
	}

	return cmd
}

func (controller *VirtualHostController) Create() *cobra.Command {
	var hostnameStr, typeStr, parentHostnameStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname": hostnameStr,
			}

			if typeStr != "" {
				requestBody["type"] = typeStr
			}

			if parentHostnameStr != "" {
				requestBody["parentHostname"] = parentHostnameStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(
		&typeStr, "type", "t", "", "VirtualHostType (top-level|subdomain|alias)",
	)
	cmd.Flags().StringVarP(&parentHostnameStr, "parent", "p", "", "ParentHostname")
	return cmd
}

func (controller *VirtualHostController) Delete() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname": hostnameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	return cmd
}

func (controller *VirtualHostController) ReadWithMappings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadVirtualHostsWithMappings",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.ReadWithMappings(),
			)
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
			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)

			err := useCase.CreateMapping(
				mappingQueryRepo,
				mappingCmdRepo,
				vhostQueryRepo,
				servicesQueryRepo,
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
