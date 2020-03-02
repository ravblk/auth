package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//RootCmd
var RootCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth - service registration, authorization, authetentication",
	Long:  `Please select command service or migrate`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
