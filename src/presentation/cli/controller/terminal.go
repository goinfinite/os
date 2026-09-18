package cliController

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	terminalSessionInfra "github.com/goinfinite/os/src/infra/terminalSession"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type TerminalController struct {
	terminalSessionLiaison *liaison.TerminalSessionLiaison
}

func NewTerminalController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *TerminalController {
	return &TerminalController{
		terminalSessionLiaison: liaison.NewTerminalSessionLiaison(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (controller *TerminalController) Read() *cobra.Command {
	var idStr, accountIdStr string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "ReadTerminalSessions",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]any{}
			if idStr != "" {
				requestBody["id"] = idStr
			}
			if accountIdStr != "" {
				requestBody["accountId"] = accountIdStr
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.terminalSessionLiaison.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "TerminalSessionId")
	cmd.Flags().StringVarP(&accountIdStr, "account-id", "a", "", "AccountId")
	return cmd
}

func (controller *TerminalController) Create() *cobra.Command {
	var accountIdStr, workingDirStr, commandStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateTerminalSession",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]any{
				"accountId": accountIdStr,
			}
			if workingDirStr != "" {
				requestBody["workingDir"] = workingDirStr
			}
			if commandStr != "" {
				requestBody["command"] = commandStr
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.terminalSessionLiaison.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "a", "", "AccountId")
	cliHelper.MarkRequiredFlag(cmd, "account-id")
	cmd.Flags().StringVarP(&workingDirStr, "working-dir", "w", "", "WorkingDirectory")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	return cmd
}

func (controller *TerminalController) Delete() *cobra.Command {
	var idStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteTerminalSession",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]any{"id": idStr}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.terminalSessionLiaison.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "TerminalSessionId")
	cliHelper.MarkRequiredFlag(cmd, "id")
	return cmd
}

func (controller *TerminalController) Attach() *cobra.Command {
	var idStr string

	cmd := &cobra.Command{
		Use:   "attach",
		Short: "AttachTerminalSession",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]any{"id": idStr}
			responseOutput := controller.terminalSessionLiaison.ReadFirst(requestBody)
			if responseOutput.Status != tkPresentation.LiaisonResponseStatusSuccess {
				tkPresentation.LiaisonCliResponseRenderer(responseOutput)
				return
			}

			terminalSessionEntity, assertOk := responseOutput.Body.(entity.TerminalSession)
			if !assertOk {
				fmt.Println("TerminalSessionNotFound")
				os.Exit(1)
			}

			tmuxSessionName := terminalSessionInfra.TerminalSessionNamePrefix +
				terminalSessionEntity.Id.String()
			attachCommand := "tmux attach -t " + tmuxSessionName

			suBinaryPath, err := exec.LookPath("su")
			if err != nil {
				fmt.Println("FindSuBinaryError: ", err)
				os.Exit(1)
			}

			execErr := syscall.Exec(
				suBinaryPath,
				[]string{
					"su", "-", terminalSessionEntity.AccountUsername.String(),
					"-c", attachCommand,
				},
				os.Environ(),
			)
			if execErr != nil {
				fmt.Println("AttachTerminalSessionError: ", execErr)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "TerminalSessionId")
	cliHelper.MarkRequiredFlag(cmd, "id")
	return cmd
}
