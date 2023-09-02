package database

import "os"

const (
	METADATA_DIR = ".gity"

	defaultPermission = os.FileMode(0744)

	DATABASE_DIR = "objects"

	EXECUTABLE_MODE = "100755"
	REGULAR_MODE    = "100644"
	DIRECTORY_MODE  = "40000"

	entryFormat = "Z*H40"

	TREE_TYPE  = "tree"
	ENTRY_TYPE = "entry"
)
