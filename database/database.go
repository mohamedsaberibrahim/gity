package database

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Database struct {
	path string
}

func (db *Database) New(path string) {
	db.path = path
}

func (db *Database) Store(object ObjectInterface) error {
	var encoded_data bytes.Buffer

	data := []byte(fmt.Sprintf("%s %d\x00", object.GetType(), len(object.ToString())))
	data = append(data, object.ToString()...)

	// Encode the data
	if err := binary.Write(&encoded_data, binary.BigEndian, data); err != nil {
		return err
	}

	objectId := sha1.Sum(encoded_data.Bytes())
	object.SetOid(objectId[:])
	return db.write_object(objectId[:], encoded_data.Bytes())
}

func (db *Database) write_object(oid []byte, content []byte) error {
	oid_hex := fmt.Sprintf("%x", oid)

	object_dir := strings.Join([]string{db.path, oid_hex[0:2]}, string(os.PathSeparator))
	object_path := strings.Join([]string{object_dir, oid_hex[2:]}, string(os.PathSeparator))

	if err := os.MkdirAll(object_dir, defaultPermission); err != nil {
		return err
	}

	tmp_dir := os.TempDir()
	tmp_file, err := os.CreateTemp(tmp_dir, db.generateTmpObjectName(oid_hex))
	if err != nil {
		return err
	}
	var compressed_data bytes.Buffer

	zlibWriter := zlib.NewWriter(&compressed_data)
	if _, err := zlibWriter.Write(content); err != nil {
		return err
	}
	zlibWriter.Close()

	if _, err := tmp_file.Write(compressed_data.Bytes()); err != nil {
		return err
	}

	tmp_file.Close()

	if err := os.Rename(tmp_file.Name(), object_path); err != nil {
		return err
	}
	return nil
}

func (db *Database) generateTmpObjectName(oid string) string {
	return fmt.Sprintf("tmp_object_%x", oid[0:5])
}
