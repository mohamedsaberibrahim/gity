package commandOps

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mohamedsaberibrahim/gity/app"
	"github.com/mohamedsaberibrahim/gity/database"
)

type Commit struct {
	repo app.Repository
}

func (c Commit) Run(args []string, message string) {
	fmt.Println("commit called")

	dir, err := os.Getwd()
	git_path := strings.Join([]string{dir, database.METADATA_DIR}, string(os.PathSeparator))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
	}

	c.repo = app.Repository{}
	c.repo.New(git_path)

	c.repo.Index.Load()
	_, entries := c.repo.Index.GetSortedEntries()
	root := database.Tree{}.Build(entries)
	root.Traverse(c.repo.Database.Store)
	author := database.Author{}
	author.New(os.Getenv("GIT_AUTHOR_NAME"), os.Getenv("GIT_AUTHOR_EMAIL"), time.Now())

	parent, _ := c.repo.Refs.ReadHead()
	commit := database.Commit{}
	commit.New(parent, root.GetOid(), author, message)
	c.repo.Database.Store(&commit)

	err = c.repo.Refs.UpdateHead(commit.GetOid())
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
