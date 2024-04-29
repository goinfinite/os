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

func (controller VirtualHostController) GetVirtualHosts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetVirtualHosts",
		Run: func(cmd *cobra.Command, args []string) {
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostsList, err := useCase.GetVirtualHosts(vhostQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsList)
		},
	}

	return cmd
}

func (controller VirtualHostController) CreateVirtualHost() *cobra.Command {
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

			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

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

func (controller VirtualHostController) DeleteVirtualHost() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)

			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

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

func (controller VirtualHostController) GetVirtualHostsWithMappings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetVirtualHostsWithMappings",
		Run: func(cmd *cobra.Command, args []string) {
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostsList, err := useCase.GetVirtualHostsWithMappings(vhostQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsList)
		},
	}

	return cmd
}

func (controller VirtualHostController) CreateVirtualHostMapping() *cobra.Command {
	var hostnameStr string
	var pathStr string
	var matchPatternStr string
	var targetTypeStr string
	var targetServiceStr string
	var targetUrlStr string
	var targetHttpResponseCode uint
	var targetInlineHtmlContent string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMapping",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			path := valueObject.NewMappingPathPanic(pathStr)
			targetType := valueObject.NewMappingTargetTypePanic(targetTypeStr)

			matchPattern := valueObject.NewMappingMatchPatternPanic("begins-with")
			if matchPatternStr != "" {
				matchPattern = valueObject.NewMappingMatchPatternPanic(matchPatternStr)
			}

			var targetServicePtr *valueObject.ServiceName
			if targetServiceStr != "" {
				targetService := valueObject.NewServiceNamePanic(targetServiceStr)
				targetServicePtr = &targetService
			}

			var targetUrlPtr *valueObject.Url
			if targetUrlStr != "" {
				targetUrl := valueObject.NewUrlPanic(targetUrlStr)
				targetUrlPtr = &targetUrl
			}

			var targetHttpResponseCodePtr *valueObject.HttpResponseCode
			if targetHttpResponseCode != 0 {
				targetHttpResponseCode := valueObject.NewHttpResponseCodePanic(
					targetHttpResponseCode,
				)
				targetHttpResponseCodePtr = &targetHttpResponseCode
			}

			var targetInlineHtmlContentPtr *valueObject.InlineHtmlContent
			if targetInlineHtmlContent != "" {
				targetInlineHtmlContent := valueObject.NewInlineHtmlContentPanic(
					targetInlineHtmlContent,
				)
				targetInlineHtmlContentPtr = &targetInlineHtmlContent
			}

			createMappingDto := dto.NewCreateMapping(
				hostname,
				path,
				matchPattern,
				targetType,
				targetServicePtr,
				targetUrlPtr,
				targetHttpResponseCodePtr,
				targetInlineHtmlContentPtr,
			)

			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
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
	cmd.Flags().StringVarP(&matchPatternStr, "match", "m", "", "MatchPattern (begins-with|contains|ends-with)")
	cmd.Flags().StringVarP(
		&targetTypeStr, "type", "t", "", "MappingTargetType (service|url|response-code)",
	)
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(
		&targetServiceStr, "service", "s", "", "TargetServiceName",
	)
	cmd.Flags().StringVarP(
		&targetUrlStr, "url", "u", "", "TargetUrl",
	)
	cmd.Flags().UintVarP(
		&targetHttpResponseCode, "response-code", "r", 0, "TargetHttpResponseCode",
	)
	cmd.Flags().StringVarP(
		&targetInlineHtmlContent, "html", "h", "", "TargetInlineHtmlContent",
	)
	return cmd
}

func (controller VirtualHostController) DeleteVirtualHostMapping() *cobra.Command {
	var hostnameStr string
	var mappingIdUint uint

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMapping",
		Run: func(cmd *cobra.Command, args []string) {
			hostname := valueObject.NewFqdnPanic(hostnameStr)
			mappingId := valueObject.NewMappingIdPanic(mappingIdUint)

			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

			err := useCase.DeleteMapping(
				vhostQueryRepo,
				vhostCmdRepo,
				hostname,
				mappingId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MappingDeleted")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHost Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().UintVarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}
