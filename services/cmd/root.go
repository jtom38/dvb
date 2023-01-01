package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "dvb",
		Short: "Docker Volume Backup to keep data safe",
		Long:  "This tool will backup your docker instances and migrate the data to a safe location for you.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("version: 0.0.2")
		},
	}
	Daemon bool
	Debug bool
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
	

	root.PersistentFlags().BoolVar(&Daemon, "daemon", false, "When True the app will stay live and not close after the job finishes.")
	root.PersistentFlags().BoolVar(&Debug, "debug", false, "Defines if the app makes any changes")
}
