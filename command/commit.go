/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mohamedsaberibrahim/gity/app"
	"github.com/mohamedsaberibrahim/gity/database"
	"github.com/spf13/cobra"
)

var message string

func init() {
	commitCmd := cobra.Command{
		Use:   "commit",
		Short: COMMIT_COMMIT_SHORT_DESC,
		Long:  COMMIT_COMMIT_LONG_DESC,
		Run:   RunCommit,
	}

	rootCmd.AddCommand(&commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given <msg> as the commit message.")
}

func RunCommit(cmd *cobra.Command, args []string) {
	fmt.Println("commit called")

	dir, err := os.Getwd()
	git_path := strings.Join([]string{dir, database.METADATA_DIR}, string(os.PathSeparator))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
	}
	repo := app.Repository{}
	repo.New(git_path)

	repo.Index.Load()
	_, entries := repo.Index.GetSortedEntries()
	root := database.Tree{}.Build(entries)
	root.Traverse(repo.Database.Store)
	author := database.Author{}
	author.New(os.Getenv("GIT_AUTHOR_NAME"), os.Getenv("GIT_AUTHOR_EMAIL"), time.Now())

	parent, _ := repo.Refs.ReadHead()
	commit := database.Commit{}
	commit.New(parent, root.GetOid(), author, message)
	repo.Database.Store(&commit)

	err = repo.Refs.UpdateHead(commit.GetOid())
	if err != nil {
		fmt.Print(err)
	}
	root_commit := ""
	if parent == nil {
		root_commit = "(root-commit) "
	}
	fmt.Printf("[%s%x] %s\n", root_commit, commit.GetOid(), commit.GetMessage())
	os.Exit(0)
}
