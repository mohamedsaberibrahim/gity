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

var message string

func init() {
	commitCmd := cobra.Command{
		Use:   "commit",
		Short: COMMIT_COMMIT_SHORT_DESC,
		Long:  COMMIT_COMMIT_LONG_DESC,
		Run: func(cmd *cobra.Command, args []string) {
			commit := commandOps.Commit{}

			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
			}
			commit.New(dir, os.Stdout, os.Stderr, os.Getenv, args, message)
			status := commit.Run()
			os.Exit(status)
		},
	}

	rootCmd.AddCommand(&commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given <msg> as the commit message.")
}
