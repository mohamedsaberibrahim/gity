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

	info, err := w.StatFile(dir)
	if err != nil {
		fmt.Println("getting stat failed: ", err)
		return []string{}
	}

	if w.isIgnoreFile(info.Name()) {
		return []string{}
	}

	var file_names []string

	if info.IsDir() {
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return []string{}
		}
		for _, file := range files {
			file_path := strings.Join([]string{dir, file.Name()}, string(os.PathSeparator))
			file_names = append(file_names, w.ListFiles(file_path)...)
		}
	} else {
		file_names = append(file_names, w.relative_path_name(dir, info.Name()))
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

func (w Workspace) isIgnoreFile(file_name string) bool {
	IGNORE_FILES := [4]string{".", "..", ".gity", ".git"}
	for _, ignore_file := range IGNORE_FILES {
		if file_name == ignore_file {
			return true
		}
	}
	return false
}

func (w Workspace) relative_path_name(dir string, file_name string) string {
	if w.root_path == dir {
		return file_name
	} else {
		relative_dir := strings.TrimPrefix(dir, w.root_path+string(os.PathSeparator))
		return relative_dir
	}
}
