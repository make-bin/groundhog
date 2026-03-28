// @AI_GENERATED
package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates the "openclaw" root command and registers subcommands.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openclaw",
		Short: "OpenClaw backend server",
	}

	cmd.AddCommand(NewGatewayCommand())
	cmd.AddCommand(NewDoctorCommand())
	cmd.AddCommand(NewOnboardCommand())
	cmd.AddCommand(NewStatusCommand())
	cmd.AddCommand(NewConfigCommand())
	cmd.AddCommand(NewCronCommand())

	return cmd
}

// @AI_GENERATED: end
