package database

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type Workspace struct {
	root_path string
}

func (w *Workspace) New(root_path string) {
	w.root_path = root_path
}

func (w *Workspace) ListFiles() []fs.DirEntry {
	files, err := os.ReadDir(w.root_path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
	}
	return files
}

func (w *Workspace) ReadFile(file fs.DirEntry) ([]byte, error) {
	full_dir := strings.Join([]string{w.root_path, file.Name()}, string(os.PathSeparator))
	return os.ReadFile(full_dir)
}
