package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

type Lockfile struct {
	file_path string
	lock_path string
	lock      *os.File
}

func (l *Lockfile) New(path string) {
	l.file_path = path
	l.lock_path = change_extension(l.file_path, ".lock")
}

func (l *Lockfile) HoldForUpdate() error {
	if l.lock != nil {
		return nil
	}
	var err error
	l.lock, err = os.OpenFile(l.lock_path, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.FileMode(0777))
	if err != nil {
		return fmt.Errorf("[Lockfile][HoldForUpdate] error: %s", err)
	}
	return nil
}

func (l *Lockfile) Write(content string) error {
	l.raise_on_stale_lock()
	_, err := l.lock.WriteString(content)
	if err != nil {
		return fmt.Errorf("[Lockfile][Write] error: %s", err)
	}
	return nil
}

func (l *Lockfile) Commit() {
	l.raise_on_stale_lock()
	l.lock.Close()
	os.Rename(l.lock_path, l.file_path)
	l.lock = nil
}

func (l *Lockfile) raise_on_stale_lock() error {
	if l.lock != nil {
		return fmt.Errorf("[Lockfile][raise_on_stale_lock] path: %s", l.lock_path)
	}
	return nil
}

func change_extension(filePath string, newExtension string) string {
	fileExt := filepath.Ext(filePath)
	newFilePath := filePath[:len(filePath)-len(fileExt)] + newExtension
	return newFilePath
}
