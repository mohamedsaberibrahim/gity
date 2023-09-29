package commandOps

import (
	"fmt"
	"os"
)

type Init struct {
}

func (i Init) Run(args []string) {
	fmt.Println("init called")
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	metaDataDir := path + "/" + ".gity"

	if err := os.Mkdir(metaDataDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
		os.Exit(1)
	}

	if err := os.Mkdir(metaDataDir+`/objects`, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	}

	if err := os.Mkdir(metaDataDir+`/refs`, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	}
	fmt.Printf("Initialized empty gity repository in %s\n", path)
}
