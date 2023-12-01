package cliController

import (
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetFilesController() *cobra.Command {
	var unixFilePath string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetFiles",
		Run: func(cmd *cobra.Command, args []string) {
			unixFilePath := valueObject.NewUnixFilePathPanic(unixFilePath)

			filesQueryRepo := infra.FilesQueryRepo{}
			filesList, err := useCase.GetFiles(filesQueryRepo, unixFilePath)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, filesList)
		},
	}

	cmd.Flags().StringVarP(&unixFilePath, "filePath", "f", "", "UnixFilePath")
	cmd.MarkFlagRequired("filePath")

	return cmd
}
