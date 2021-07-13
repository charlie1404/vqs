package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of cli",
		Long:  "Print the version of cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", CURRENT_VERSION)
		},
	})
}
