/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"github.com/mohamedsaberibrahim/gity/commandOps"
	"github.com/spf13/cobra"
)

func init() {
	addCmd := cobra.Command{
		Use:   "add",
		Short: ADD_COMMIT_SHORT_DESC,
		Long:  ADD_COMMIT_LONG_DESC,
		Run: func(cmd *cobra.Command, args []string) {
			command := commandOps.Add{}
			command.Run(args)
		},
	}
	rootCmd.AddCommand(&addCmd)
}
