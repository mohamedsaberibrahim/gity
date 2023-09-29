/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"fmt"
	"os"

	"github.com/mohamedsaberibrahim/gity/commandOps"
	"github.com/spf13/cobra"
)

func init() {
	addCmd := cobra.Command{
		Use:   "add",
		Short: ADD_COMMIT_SHORT_DESC,
		Long:  ADD_COMMIT_LONG_DESC,
		Run: func(cmd *cobra.Command, args []string) {
			add := commandOps.Add{}

			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
			}
			add.New(dir, os.Stdout, os.Stderr, args)
			status := add.Run()
			os.Exit(status)
		},
	}
	rootCmd.AddCommand(&addCmd)
}
