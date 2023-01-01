package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts running the application.",
		Long:  "Starts running the application.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ConfigPath: %v\n", ConfigPath)
			fmt.Printf("Daemon: %v\n", Daemon)

			info, err := os.Stat(ConfigPath)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
			fmt.Printf("info: %v\n", info)
		},
	}
	ConfigPath string
)


func init() {
	startCmd.Flags().StringVar(&ConfigPath, "config-path", "", "Defines what config file should be loaded")
}

