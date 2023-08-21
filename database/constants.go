package database

import "os"

const (
	metadataDir = ".gity"

	defaultPermission = os.FileMode(0744)

	databaseDir = "objects"
)
