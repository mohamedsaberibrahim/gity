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

func (w *Workspace) ListFiles(dir string) []string {
	if dir == "" {
		dir = w.root_path
	}
	// fmt.Print("current directory: ", dir, "\n")
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read the current directory - %v\n", err)
	}
	var file_names []string
	for _, file := range files {
		if isIgnoreFile(file) {
			continue
		}
		if file.Type().IsRegular() {
			relative_path_name := w.relative_path_name(dir, file)
			file_names = append(file_names, relative_path_name)
		} else if file.IsDir() {
			file_path := strings.Join([]string{dir, file.Name()}, string(os.PathSeparator))
			file_names = append(file_names, w.ListFiles(file_path)...)
		}
	}
	return file_names
}

func (w *Workspace) ReadFile(file string) ([]byte, error) {
	full_dir := strings.Join([]string{w.root_path, file}, string(os.PathSeparator))
	return os.ReadFile(full_dir)
}

func (w *Workspace) StatFile(file string) (fs.FileInfo, error) {
	return os.Stat(file)
}

func isIgnoreFile(file fs.DirEntry) bool {
	IGNORE_FILES := [4]string{".", "..", ".gity", ".git"}
	for _, ignore_file := range IGNORE_FILES {
		if file.Name() == ignore_file {
			return true
		}
	}
	return false
}

func (w *Workspace) relative_path_name(dir string, file fs.DirEntry) string {
	if w.root_path == dir {
		return file.Name()
	} else {
		relative_dir := strings.TrimPrefix(dir, w.root_path+string(os.PathSeparator))
		return strings.Join([]string{relative_dir, file.Name()}, string(os.PathSeparator))
	}
}
