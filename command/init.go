/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	initCmd := cobra.Command{
		Use:   "init",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: RunInit,
	}
	rootCmd.AddCommand(&initCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	fmt.Println("init called")
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	metaDataDir := path + "/" + ".gity"

	if err := os.Mkdir(metaDataDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
		os.Exit(1)
	}

	if err := os.Mkdir(metaDataDir+`/objects`, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	}

	if err := os.Mkdir(metaDataDir+`/refs`, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	}
	fmt.Printf("Initialized empty gity repository in %s\n", path)
}
