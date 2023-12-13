package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetVirtualHostsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetVirtualHosts",
		Run: func(cmd *cobra.Command, args []string) {
			vhostQueryRepo := infra.VirtualHostQueryRepo{}
			vhostsList, err := useCase.GetVirtualHosts(vhostQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsList)
		},
	}

	return cmd
}

func AddVirtualHostController() *cobra.Command {
	var hostnameStr string
	var typeStr string
	var parentHostnameStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddVirtualHost",
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

			addVirtualHostDto := dto.NewAddVirtualHost(
				hostname,
				vhostType,
				parentHostnamePtr,
			)

			vhostQueryRepo := infra.VirtualHostQueryRepo{}
			vhostCmdRepo := infra.VirtualHostCmdRepo{}

			err := useCase.AddVirtualHost(
				vhostQueryRepo,
				vhostCmdRepo,
				addVirtualHostDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "VirtualHostAdded")
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

func GetVirtualHostsWithMappingsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetVirtualHostsWithMappings",
		Run: func(cmd *cobra.Command, args []string) {
			vhostQueryRepo := infra.VirtualHostQueryRepo{}
			vhostsList, err := useCase.GetVirtualHostsWithMappings(vhostQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, vhostsList)
		},
	}

	return cmd
}
