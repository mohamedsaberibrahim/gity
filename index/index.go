package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"sort"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/helper"
)

type Index struct {
	entries      map[string]Entry
	lockfile     helper.Lockfile
	keys         []string
	encoded_data bytes.Buffer
}

func (i *Index) New(path string) {
	i.entries = map[string]Entry{}
	i.keys = []string{}
	i.lockfile.New(path)
}

func (i *Index) Add(path_name string, oid []byte, stat syscall.Stat_t) {
	fmt.Println("Adding new entry to index: ", path_name)
	entry := Entry{}
	entry.New(path_name, oid, stat)
	i.entries[path_name] = entry
	i.keys = append(i.keys, path_name)
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

	sort.Strings(i.keys)

	for _, key := range i.keys {
		entry := i.entries[key]
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
	fmt.Println("Before writing into data: ", i.encoded_data)
	if err := binary.Write(&i.encoded_data, binary.BigEndian, []byte(data)); err != nil {
		fmt.Printf("Error: failed to write data - %v\n", err)
		return err
	}
	fmt.Println("After writing into data: ", i.encoded_data)
	i.lockfile.Write(data)
	return nil
}

func (i *Index) finish_write() {
	objectId := sha1.Sum(i.encoded_data.Bytes())
	oid_hex := fmt.Sprintf("%x", objectId[:])
	fmt.Println("objectId: ", oid_hex)
	i.lockfile.Write(string(objectId[:]))
	i.lockfile.Commit()
}
