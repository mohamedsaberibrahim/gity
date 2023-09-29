package commandOps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/app"
	"github.com/mohamedsaberibrahim/gity/database"
	"github.com/mohamedsaberibrahim/gity/index"
)

type Add struct {
	repo app.Repository
}

func (a Add) Run(args []string) {
	fmt.Println("add called", args)
	dir, err := os.Getwd()
	git_path := strings.Join([]string{dir, database.METADATA_DIR}, string(os.PathSeparator))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
	}

	a.repo = app.Repository{}
	a.repo.New(git_path)

	_, err = a.repo.Index.LoadForUpdate()
	if err != nil {
		fmt.Printf("fatal: %s\n\nAnother jit process seems to be running in this repository.\nPlease make sure all processes are terminated then try again.\nIf it still fails, a jit process may have crashed in this\nrepository earlier: remove the file manually to continue.\n", err)
		os.Exit(128)
	}

	paths, err := get_paths(a.repo.Workspace, args)
	if err != nil {
		fmt.Printf("fatal: %s\n", err)
		a.repo.Index.ReleaseLock()
		os.Exit(128)
	}
	err = add_entries(paths, &a.repo.Workspace, &a.repo.Database, &a.repo.Index)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		fmt.Println("fatal: adding files failed")
		a.repo.Index.ReleaseLock()
		os.Exit(128)
	}
	a.repo.Index.WriteUpdates()
	os.Exit(0)
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

func add_entries(paths []string, workspace *database.Workspace, db *database.Database, index *index.Index) error {
	fmt.Println("Adding entries: ", paths)
	for _, file_path := range paths {
		var st syscall.Stat_t
		if err := syscall.Stat(file_path, &st); err != nil {
			log.Fatal(err)
		}

		data, err := (*workspace).ReadFile(file_path)
		if err != nil {
			return fmt.Errorf("%s\nerror: unable to index file %s", err, file_path)
		}

		blob := database.Blob{}
		blob.New(data)
		(*db).Store(&blob)
		(*index).Add(file_path, blob.GetOid(), st)
		fmt.Println("After adding to index: ", file_path)
	}
	return nil
}
