package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigPath     string
	Daemon         bool
	Debug          bool
	Version        string = "0.0.7"

	root = &cobra.Command{
		Use:   "dvb",
		Short: "Docker Volume Backup to keep data safe",
		Long:  "This tool will backup your docker instances and migrate the data to a safe location for you.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("use --help for details on this application.")
		},
	}
)

func Execute() {
	err := root.Execute()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(startCmd)
	root.AddCommand(versionCmd)
	//root.AddCommand(installCmd)

	//root.PersistentFlags().BoolVar(&Daemon, "daemon", false, "When True the app will stay live and not close after the job finishes.")
	//root.PersistentFlags().BoolVar(&Debug, "debug", false, "Defines if the app makes any changes")
}
