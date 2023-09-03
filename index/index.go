package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/helper"
)

type Index struct {
	entries      map[string]Entry
	lockfile     helper.Lockfile
	encoded_data bytes.Buffer
}

func (i *Index) New(path string) {
	i.entries = map[string]Entry{}
	i.lockfile.New(path)
}

func (i *Index) Add(path_name string, oid []byte, stat syscall.Stat_t) {
	entry := Entry{}
	entry.New(path_name, oid, stat)
	i.entries[path_name] = entry
}

func (i *Index) WriteUpdates() bool {
	if ok, _ := i.lockfile.HoldForUpdate(); !ok {
		fmt.Println("Could not lock index file")
		return false
	}
	i.begin_write()
	header := make([]byte, 12)
	copy(header[:4], "DIRC")
	binary.BigEndian.PutUint32(header[4:8], 2)
	binary.BigEndian.PutUint32(header[8:12], uint32(len(i.entries)))
	fmt.Println("header: ", string(header))
	i.write(string(header))

	for _, entry := range i.entries {
		fmt.Println(entry.ToString())
		i.write(entry.ToString())
	}
	i.finish_write()
	return true
}

func (i *Index) begin_write() {
	i.encoded_data.Reset()
}

func (i *Index) write(data string) error {
	// if err := binary.Write(&i.encoded_data, binary.BigEndian, data); err != nil {
	// 	return err
	// }
	i.lockfile.Write(data)
	return nil
}

func (i *Index) finish_write() {
	objectId := sha1.Sum(i.encoded_data.Bytes())
	fmt.Println("objectId: ", string(objectId[:]))
	oid_hex := fmt.Sprintf("%x", objectId[:])
	fmt.Println("objectId: ", oid_hex)
	i.lockfile.Write(string(objectId[:]))
	i.lockfile.Commit()
}
