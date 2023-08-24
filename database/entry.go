package database

import (
	"fmt"
	"io/fs"
)

const (
	EXECUTABLE_MODE = "100755"
	REGULAR_MODE    = "100644"
)

type Entry struct {
	name string
	oid  []byte
	stat fs.FileInfo
}

func (e *Entry) New(name string, oid []byte, stat fs.FileInfo) {
	e.name = name
	e.oid = oid
	e.stat = stat
}

func (e *Entry) GetName() string {
	return e.name
}

func (e *Entry) GetOid() []byte {
	return e.oid
}

func (e *Entry) GetMode() string {
	mode := REGULAR_MODE
	if fmt.Sprintf("10%04o", e.stat.Mode().Perm()) == EXECUTABLE_MODE {
		mode = EXECUTABLE_MODE
	}
	return mode
}
