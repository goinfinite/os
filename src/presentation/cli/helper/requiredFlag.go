package cliHelper

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func MarkRequiredFlag(command *cobra.Command, flagName string) {
	err := command.MarkFlagRequired(flagName)
	if err != nil {
		slog.Error(
			"MarkFlagRequiredError",
			slog.String("flagName", flagName),
			slog.String("err", err.Error()),
		)
	}
}
