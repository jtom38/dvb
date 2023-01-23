package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jtom38/dvb/services/proc"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts running the application.",
		Long:  "Starts running the application.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %v\n", Version)
			fmt.Printf("ConfigPath: %v\n", ConfigPath)
			//fmt.Printf("Daemon: %v\n", Daemon)

			_, err := os.Stat(ConfigPath)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}

			client := proc.NewStartBackupClient(proc.StartBackupParams{
				ConfigPath: ConfigPath,
				Daemon:     Daemon,
			})
			err = client.RunProcess()
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	startCmd.Flags().StringVar(&ConfigPath, "config-path", "", "Defines what config file should be loaded")
}
