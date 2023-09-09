package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedsaberibrahim/gity/database"
	"github.com/mohamedsaberibrahim/gity/index"
	"github.com/mohamedsaberibrahim/gity/reference"
)

type Repository struct {
	git_path  string
	Workspace database.Workspace
	Database  database.Database
	Index     index.Index
	Refs      reference.Ref
}

func (r *Repository) New(git_path string) {
	dir := filepath.Dir(git_path)
	r.Workspace = database.Workspace{}
	r.Workspace.New(dir)
	r.Database = database.Database{}
	r.Database.New(strings.Join([]string{git_path, database.DATABASE_DIR}, string(os.PathSeparator)))
	r.Index = index.Index{}
	r.Index.New(strings.Join([]string{git_path, "index"}, string(os.PathSeparator)))
	r.Refs = reference.Ref{}
	r.Refs.New(git_path)
}
