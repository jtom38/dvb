package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays the version of the tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
)
