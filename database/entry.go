package database

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type Entry struct {
	name        string
	oid         []byte
	stat        fs.FileMode
	parent_dirs []string
}

func (e *Entry) New(name string, oid []byte, stat fs.FileMode) {
	e.name = name
	e.oid = oid
	e.stat = stat
}

func (e *Entry) ToString() string {
	return e.name
}

func (e *Entry) GetType() string {
	return ENTRY_TYPE
}

func (e *Entry) GetOid() []byte {
	return e.oid
}

func (e *Entry) GetMode() string {
	mode := REGULAR_MODE
	if fmt.Sprintf("10%04o", e.stat.Perm()) == EXECUTABLE_MODE {
		mode = EXECUTABLE_MODE
	}
	return mode
}

func (e *Entry) SetOid(oid []byte) {
	e.oid = oid
}

func (e *Entry) GetParentDirectories(path string) []string {
	if path == "" {
		path = e.name
		fmt.Println("Executing descend for path: ", path)
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return e.parent_dirs
	}
	e.GetParentDirectories(dir)
	e.parent_dirs = append(e.parent_dirs, dir)
	return e.parent_dirs
}

func GetBaseName(path string) string {
	return filepath.Base(path)
}
