package commandOps

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mohamedsaberibrahim/gity/app"
	"github.com/mohamedsaberibrahim/gity/database"
)

type Commit struct {
	repo    app.Repository
	dir     string
	stdout  io.Writer
	stderr  io.Writer
	Getenv  func(key string) string
	args    []string
	message string
}

func (c *Commit) New(dir string, stdout io.Writer, stderr io.Writer, Getenv func(key string) string, args []string, message string) {
	c.dir = dir
	c.stdout = stdout
	c.stderr = stderr
	c.Getenv = Getenv
	c.args = args
	c.message = message
}

func (c *Commit) Run() int {
	fmt.Println("commit called")
	git_path := strings.Join([]string{c.dir, database.METADATA_DIR}, string(os.PathSeparator))

	c.repo = app.Repository{}
	c.repo.New(git_path)

	c.repo.Index.Load()
	_, entries := c.repo.Index.GetSortedEntries()
	root := database.Tree{}.Build(entries)
	root.Traverse(c.repo.Database.Store)
	author := database.Author{}
	author.New(c.Getenv("GIT_AUTHOR_NAME"), c.Getenv("GIT_AUTHOR_EMAIL"), time.Now())

	parent, _ := c.repo.Refs.ReadHead()
	commit := database.Commit{}
	commit.New(parent, root.GetOid(), author, c.message)
	c.repo.Database.Store(&commit)

	c.repo.Refs.UpdateHead(commit.GetOid())
	root_commit := ""
	if parent == nil {
		root_commit = "(root-commit) "
	}
	fmt.Printf("[%s%x] %s\n", root_commit, commit.GetOid(), commit.GetMessage())
	return 0
}
