/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/database"
	"github.com/mohamedsaberibrahim/gity/globals"
	"github.com/mohamedsaberibrahim/gity/index"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("add called", args)
		dir, err := os.Getwd()
		git_path := strings.Join([]string{dir, database.METADATA_DIR}, string(os.PathSeparator))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
		}
		workspace := database.Workspace{}
		db := database.Database{}
		workspace.New(dir)

		db.New(strings.Join([]string{git_path, database.DATABASE_DIR}, string(os.PathSeparator)))
		index := index.Index{}
		index.New(strings.Join([]string{git_path, "index"}, string(os.PathSeparator)))
		index.LoadForUpdate()

		paths, err := get_paths(workspace, args)
		if err != nil {
			fmt.Printf("fatal: %s\n", err)
			index.ReleaseLock()
			os.Exit(128)
		}
		err = add_entries(paths, workspace, db, index)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			fmt.Println("fatal: adding files failed")
			index.ReleaseLock()
			os.Exit(128)
		}
		index.WriteUpdates()

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func get_paths(workspace database.Workspace, args []string) ([]string, error) {
	paths := []string{}
	for _, passed_path := range args {
		abs_path, err := filepath.Abs(passed_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
		}
		files_name, err := workspace.ListFiles(abs_path)
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, files_name...)
	}
	return paths, nil
}

func add_entries(paths []string, workspace database.Workspace, db database.Database, index index.Index) *globals.ValueError {
	for _, file_path := range paths {
		var st syscall.Stat_t
		if err := syscall.Stat(file_path, &st); err != nil {
			log.Fatal(err)
		}

		data, err := workspace.ReadFile(file_path)
		if err != nil {
			return globals.NewValueError(err.Value, fmt.Errorf("%s\nerror: unable to index file %s", err.Err, file_path))
		}

		blob := database.Blob{}
		blob.New(data)
		db.Store(&blob)
		index.Add(file_path, blob.GetOid(), st)
	}
	return nil
}
