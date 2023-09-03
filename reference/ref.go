package reference

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohamedsaberibrahim/gity/helper"
)

type Ref struct {
	git_path string
}

func (r *Ref) New(git_path string) {
	r.git_path = git_path
}

func (r *Ref) ReadHead() ([]byte, error) {
	b, err := os.ReadFile(r.getHeadPath())
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *Ref) UpdateHead(oid []byte) error {
	lockfile := helper.Lockfile{}
	lockfile.New(r.getHeadPath())
	_, err := lockfile.HoldForUpdate()
	if err != nil {
		return fmt.Errorf("Could not acquire lock on file: %s error: %e", r.getHeadPath(), err)
	}
	oid_hex := fmt.Sprintf("%x", oid)
	err = lockfile.Write(oid_hex)
	if err != nil {
		return fmt.Errorf("Writing: %s error: %e", r.getHeadPath(), err)
	}
	lockfile.Commit()
	return nil
}

func (r *Ref) getHeadPath() string {
	return strings.Join([]string{r.git_path, "HEAD"}, string(os.PathSeparator))
}
