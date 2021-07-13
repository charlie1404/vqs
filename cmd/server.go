package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Start Queue Service",
		Long:  "Start Queue Service",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting Server")
		},
	})
}
