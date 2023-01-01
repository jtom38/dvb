package cmd

import "github.com/spf13/cobra"

var (
	IsUser bool

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Installs the application and sets up basic things.",
		Long:  `Installs the applications and setups a basic config file.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	installCmd.Flags().StringVar(&ConfigPath, "config-path", "", "Defines where to generate a config file.")
	installCmd.Flags().BoolVar(&IsUser, "user", false, "Defines the target to be the user profile.")
}
