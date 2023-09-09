package helper

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Lockfile struct {
	file_path string
	lock_path string
	Lock      *os.File
}

func (l *Lockfile) New(path string) {
	l.file_path = path
	l.lock_path = change_extension(l.file_path, ".lock")
}

func (l *Lockfile) HoldForUpdate() (bool, error) {
	if l.Lock != nil {
		return false, nil
	}
	var err error
	l.Lock, err = os.OpenFile(l.lock_path, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.FileMode(0777))
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return false, fmt.Errorf("Unable to create %s: File exists.", l.lock_path)
		}
		return false, err
	}
	return true, nil
}

func (l *Lockfile) Write(content string) error {
	l.raise_on_stale_lock()
	_, err := l.Lock.WriteString(content)
	if err != nil {
		return fmt.Errorf("[Lockfile][Write] error: %s", err)
	}
	return nil
}

func (l *Lockfile) Commit() {
	l.raise_on_stale_lock()
	l.Lock.Close()
	os.Rename(l.lock_path, l.file_path)
	l.Lock = nil
}

func (l *Lockfile) Rollback() {
	l.raise_on_stale_lock()
	l.Lock.Close()
	os.Remove(l.lock_path)
	l.Lock = nil
}

func (l *Lockfile) raise_on_stale_lock() error {
	if l.Lock != nil {
		return fmt.Errorf("[Lockfile][raise_on_stale_lock] path: %s", l.lock_path)
	}
	return nil
}

func change_extension(filePath string, newExtension string) string {
	fileExt := filepath.Ext(filePath)
	newFilePath := filePath[:len(filePath)-len(fileExt)] + newExtension
	return newFilePath
}
