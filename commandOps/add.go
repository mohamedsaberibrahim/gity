package commandOps

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/app"
	"github.com/mohamedsaberibrahim/gity/database"
)

type Add struct {
	repo   app.Repository
	dir    string
	stdout io.Writer
	stderr io.Writer
	args   []string
}

func (a *Add) New(dir string, stdout io.Writer, stderr io.Writer, args []string) {
	a.dir = dir
	a.stdout = stdout
	a.stderr = stderr
	a.args = args
}

func (a *Add) Run() int {
	fmt.Println("add called", a.args)
	git_path := strings.Join([]string{a.dir, database.METADATA_DIR}, string(os.PathSeparator))

	a.repo = app.Repository{}
	a.repo.New(git_path)

	_, err := a.repo.Index.LoadForUpdate()
	if err != nil {
		fmt.Fprintf(a.stderr, "fatal: %s\n\nAnother jit process seems to be running in this repository.\nPlease make sure all processes are terminated then try again.\nIf it still fails, a jit process may have crashed in this\nrepository earlier: remove the file manually to continue.\n", err)
		return 128
	}

	paths, err := a.get_paths()
	if err != nil {
		fmt.Fprintf(a.stderr, "fatal: %s\n", err)
		a.repo.Index.ReleaseLock()
		return 128
	}
	err = a.add_entries(paths)
	if err != nil {
		fmt.Fprintf(a.stderr, "error: %s\n", err)
		fmt.Fprint(a.stderr, "fatal: adding files failed\n")
		a.repo.Index.ReleaseLock()
		return 128
	}
	a.repo.Index.WriteUpdates()
	return 0
}

func (a *Add) get_paths() ([]string, error) {
	paths := []string{}
	for _, passed_path := range a.args {
		abs_path, err := filepath.Abs(passed_path)
		if err != nil {
			fmt.Fprintf(a.stderr, "Error: failed to read the current directory - %v\n", err)
		}
		files_name, err := a.repo.Workspace.ListFiles(abs_path)
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, files_name...)
	}
	return paths, nil
}

func (a *Add) add_entries(paths []string) error {
	fmt.Println("Adding entries: ", paths)
	for _, file_path := range paths {
		var st syscall.Stat_t
		if err := syscall.Stat(file_path, &st); err != nil {
			log.Fatal(err)
		}

		data, err := a.repo.Workspace.ReadFile(file_path)
		if err != nil {
			return fmt.Errorf("%s\nerror: unable to index file %s", err, file_path)
		}

		blob := database.Blob{}
		blob.New(data)
		a.repo.Database.Store(&blob)
		a.repo.Index.Add(file_path, blob.GetOid(), st)
		fmt.Fprintf(a.stdout, "After adding to index: %s", file_path)
	}
	return nil
}
