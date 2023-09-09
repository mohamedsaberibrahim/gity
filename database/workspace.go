package database

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/mohamedsaberibrahim/gity/globals"
)

type Workspace struct {
	root_path string
}

func (w *Workspace) New(root_path string) {
	w.root_path = root_path
}

func (w *Workspace) ListFiles(dir string) ([]string, error) {
	if dir == "" {
		dir = w.root_path
	}

	info, err := w.StatFile(dir)
	if err != nil {
		relative_file_name := w.relative_path_name(dir, "")
		if err.Value == globals.ErrPermissionDenied {
			return []string{}, globals.NewValueError(globals.ErrPermissionDenied, fmt.Errorf("stat(%s): Permission denied", relative_file_name))
		}
		if err.Value == globals.ErrFileNotFound {
			return []string{}, globals.NewValueError(globals.ErrFileNotFound, fmt.Errorf("pathspec %s did not match any files", relative_file_name))
		}
	}

	if w.isIgnoreFile(info.Name()) {
		return []string{}, nil
	}

	var file_names []string

	if info.IsDir() {
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return []string{}, nil
		}
		for _, file := range files {
			file_path := strings.Join([]string{dir, file.Name()}, string(os.PathSeparator))
			result, err := w.ListFiles(file_path)
			if err != nil {
				return []string{}, err
			}
			file_names = append(file_names, result...)
		}
	} else if info.Mode().IsRegular() {
		file_names = append(file_names, w.relative_path_name(dir, info.Name()))
	} else {
		return []string{}, fmt.Errorf("pathspec %s did not match any files", dir)
	}

	return file_names, nil
}

func (w *Workspace) ReadFile(file string) ([]byte, error) {
	full_dir := strings.Join([]string{w.root_path, file}, string(os.PathSeparator))
	data, err := os.ReadFile(full_dir)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return nil, fmt.Errorf("open(%s): Permission denied", file)
		} else {
			return nil, err
		}
	}
	return data, nil
}

func (w *Workspace) StatFile(file string) (fs.FileInfo, *globals.ValueError) {
	stat, err := os.Stat(file)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, globals.NewValueError(globals.ErrFileNotFound, fmt.Errorf(""))
		} else if errors.Is(err, fs.ErrPermission) {
			return nil, globals.NewValueError(globals.ErrPermissionDenied, fmt.Errorf("stat(%s): Permission denied", file))
		}
	}
	return stat, nil
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
