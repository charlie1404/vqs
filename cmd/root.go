package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CURRENT_VERSION string = "0.0.0"

var rootCmd = &cobra.Command{
	Use:   "vqueue",
	Short: "VQueue is simple queue management service, based on memory mapped files.",
	Long:  "VQueue is simple queue management service, based on memory mapped files.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
