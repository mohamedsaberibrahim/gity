package index

import (
	"encoding/binary"
	"fmt"
	"math"
	"syscall"
)

const (
	REGULAR_MODE    = 0100644
	EXECUTABLE_MODE = 0100755
	MAX_PATH_SIZE   = 0xfff
	INDEX_FILE      = "index"
	ENTRY_BLOCK     = 8
)

type Entry struct {
	ctime      int64
	ctime_nsec int64
	mtime      int64
	mtime_nsec int64
	dev        uint64
	ino        uint64
	file_size  int64
	uid        uint32
	gid        uint32
	flags      uint16
	path       string
	oid        []byte
	mode       uint32
}

func (e *Entry) New(path_name string, oid []byte, stat syscall.Stat_t) {
	e.ctime = stat.Ctim.Sec
	e.ctime_nsec = int64(stat.Ctim.Nsec)
	e.mtime = stat.Mtim.Sec
	e.mtime_nsec = int64(stat.Mtim.Nsec)
	e.dev = stat.Dev
	e.ino = stat.Ino
	if stat.Mode == syscall.S_IXUSR {
		e.mode = EXECUTABLE_MODE
	} else {
		e.mode = REGULAR_MODE
	}
	e.file_size = stat.Size
	e.uid = stat.Uid
	e.gid = stat.Gid
	e.flags = uint16(math.Min(float64(stat.Size), float64(MAX_PATH_SIZE)))
	e.path = path_name
	e.oid = oid
}

func (e *Entry) ToString() string {
	data := make([]byte, 70)
	binary.BigEndian.PutUint32(data[:4], uint32(e.ctime))
	binary.BigEndian.PutUint32(data[4:8], uint32(e.ctime_nsec))
	binary.BigEndian.PutUint32(data[8:12], uint32(e.mtime))
	binary.BigEndian.PutUint32(data[12:16], uint32(e.mtime_nsec))
	binary.BigEndian.PutUint32(data[16:20], uint32(e.dev))
	binary.BigEndian.PutUint32(data[20:24], uint32(e.ino))
	binary.BigEndian.PutUint32(data[24:28], uint32(e.mode))
	binary.BigEndian.PutUint32(data[28:32], uint32(e.uid))
	binary.BigEndian.PutUint32(data[32:36], uint32(e.gid))
	binary.BigEndian.PutUint32(data[36:40], uint32(e.file_size))
	// fmt.Println("data: ", string(data))
	copy(data[40:60], e.oid)
	fmt.Println("flags: ", e.flags)
	// Writing flags to data
	binary.BigEndian.PutUint16(data[60:62], e.flags)
	// Writing path to data
	copy(data[62:], []byte(e.path))
	// Printing data path

	padLength := (ENTRY_BLOCK - len(data)%ENTRY_BLOCK) % ENTRY_BLOCK
	padding := make([]byte, padLength)
	data = append(data, padding...)
	// Printing the data
	fmt.Println("data: ", string(data))
	// Encode the data
	return string(data)
}
