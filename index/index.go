package index

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/mohamedsaberibrahim/gity/helper"
)

const (
	HEADER_SIZE    = 12
	SIGNATURE      = "DIRC"
	VERSION        = 2
	ENTRY_MIN_SIZE = 64
)

type Index struct {
	entries      map[string]Entry
	lockfile     helper.Lockfile
	path         string
	encoded_data bytes.Buffer
	changed      bool
}

func (i *Index) New(path string) {
	i.entries = map[string]Entry{}
	i.lockfile.New(path)
	i.path = path
	i.changed = false
}

func (i *Index) LoadForUpdate() bool {
	held_successfully, _ := i.lockfile.HoldForUpdate()
	if !held_successfully {
		// fmt.Println("Could not lock index file")
		return false
	}
	i.Load()
	return true
}

func (i *Index) Load() {
	i.Clear()
	reader := Checksum{}
	file := i.open_index_file()
	if file != nil {
		reader.New(file)
		entries_count := i.read_header(&reader)
		i.read_entries(&reader, entries_count)
		reader.VerifyChecksum()
		file.Close()
	}
}

func (i *Index) Clear() {
	i.entries = map[string]Entry{}
	i.changed = false
}

func (i *Index) Add(path_name string, oid []byte, stat syscall.Stat_t) {
	fmt.Println("Restoring entry to index: ", path_name)
	entry := Entry{}
	entry.New(path_name, oid, stat)
	i.store_entry(entry)
	i.changed = true
}

func (i *Index) WriteUpdates() bool {
	if !i.changed {
		i.lockfile.Rollback()
		return true
	}
	writer := Checksum{}
	writer.New(i.lockfile.Lock)
	header := make([]byte, 12)
	copy(header[:4], SIGNATURE)
	binary.BigEndian.PutUint32(header[4:8], VERSION)
	binary.BigEndian.PutUint32(header[8:12], uint32(len(i.entries)))
	writer.Write(header)

	keys := []string{}
	for key := range i.entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		entry := i.entries[key]
		fmt.Println("Writing entry: ", entry.GetPath())
		writer.Write([]byte(entry.ToString()))
	}

	writer.WriteChecksum()
	i.lockfile.Commit()
	i.changed = false
	return true
}

func (i *Index) open_index_file() *os.File {
	file, err := os.Open(i.path)
	if err != nil {
		return nil
	}
	return file
}

func (i *Index) read_header(c *Checksum) uint32 {
	data, err := c.Read(HEADER_SIZE)
	if err != nil {
		fmt.Println("Error: failed to read header - ", err)
	}
	signature := string(data[:4])
	version := binary.BigEndian.Uint32(data[4:8])
	count := binary.BigEndian.Uint32(data[8:12])

	if signature != SIGNATURE {
		fmt.Println("Error: invalid signature")
	}
	if version != VERSION {
		fmt.Println("Error: invalid version")
	}
	fmt.Println("count: ", count)
	return count
}

func (i *Index) read_entries(c *Checksum, entries_count uint32) {
	fmt.Println("Reading entries: ", entries_count)
	for j := 0; j < int(entries_count); j++ {
		data, _ := c.Read(ENTRY_MIN_SIZE)
		fmt.Println("last byte: ", data[len(data)-1])
		for data[len(data)-1] != 0 {
			extra_byte, err := c.Read(ENTRY_BLOCK)
			if err != nil {
				fmt.Println("Error: failed to read entry - ", err)
				break
			}
			data = append(data, extra_byte...)
			// printing data size
			fmt.Println("data size: ", len(data))
		}
		e := Entry{}
		entry := e.Parse(data)
		fmt.Println("Entry path: ", data, entry.GetPath())
		i.store_entry(entry)
	}
}

func (i *Index) store_entry(entry Entry) {
	fmt.Println("Entries size before: ", len(i.entries), entry.GetPath())
	i.entries[entry.path] = entry
	fmt.Println("Entries size after: ", len(i.entries), entry.GetPath())
}
