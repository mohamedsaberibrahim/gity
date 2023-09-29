package commandOps

import (
	"fmt"
	"io"
	"os"
)

type Init struct {
	dir    string
	stdout io.Writer
	stderr io.Writer
	args   []string
}

func (i *Init) New(dir string, stdout io.Writer, stderr io.Writer, args []string) {
	i.dir = dir
	i.stdout = stdout
	i.stderr = stderr
	i.args = args
}

func (i *Init) Run() int {
	fmt.Println("init called")
	metaDataDir := i.dir + "/" + ".gity"

	if err := os.Mkdir(metaDataDir, os.ModePerm); err != nil {
		fmt.Fprintf(i.stderr, "fatal: %s\n", err)
		return 1
	}

	if err := os.Mkdir(metaDataDir+`/objects`, os.ModePerm); err != nil {
		fmt.Fprintf(i.stderr, "fatal: %s\n", err)
		return 1
	}

	if err := os.Mkdir(metaDataDir+`/refs`, os.ModePerm); err != nil {
		fmt.Fprintf(i.stderr, "fatal: %s\n", err)
		return 1
	}
	fmt.Fprintf(i.stdout, "Initialized empty gity repository in %s\n", i.dir)
	return 0
}
