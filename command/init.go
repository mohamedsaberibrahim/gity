/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"os"

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
			dir, err := os.Getwd()
			if err != nil {
				os.Exit(1)
			}
			init.New(dir, os.Stdout, os.Stderr, args)
			status := init.Run()
			os.Exit(status)
		},
	}
	rootCmd.AddCommand(&initCmd)
}
