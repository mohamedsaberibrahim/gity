/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
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
			commit.Run(args, message)
		},
	}

	rootCmd.AddCommand(&commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given <msg> as the commit message.")
}
