/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"github.com/mohamedsaberibrahim/gity/commandOps"
	"github.com/spf13/cobra"
)

func init() {
	initCmd := cobra.Command{
		Use:   "init",
		Short: INIT_COMMIT_SHORT_DESC,
		Long:  INIT_COMMIT_LONG_DESC,
		Run: func(cmd *cobra.Command, args []string) {
			init := commandOps.Init{}
			init.Run(args)
		},
	}
	rootCmd.AddCommand(&initCmd)
}
