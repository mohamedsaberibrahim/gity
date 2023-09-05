package index

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
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
	dev        uint32
	ino        uint32
	file_size  uint32
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
	e.dev = uint32(stat.Dev)
	e.ino = uint32(stat.Ino)
	if stat.Mode == syscall.S_IXUSR {
		e.mode = EXECUTABLE_MODE
	} else {
		e.mode = REGULAR_MODE
	}
	e.file_size = uint32(stat.Size)
	e.uid = stat.Uid
	e.gid = stat.Gid
	e.flags = uint16(math.Min(float64(len([]byte(path_name))), float64(MAX_PATH_SIZE)))
	e.path = path_name
	e.oid = oid
}

func (e *Entry) ToString() string {
	data := make([]byte, 62+len([]byte(e.path)))
	binary.BigEndian.PutUint32(data[:4], uint32(e.ctime))
	binary.BigEndian.PutUint32(data[4:8], uint32(e.ctime_nsec))
	binary.BigEndian.PutUint32(data[8:12], uint32(e.mtime))
	binary.BigEndian.PutUint32(data[12:16], uint32(e.mtime_nsec))
	binary.BigEndian.PutUint32(data[16:20], e.dev)
	binary.BigEndian.PutUint32(data[20:24], e.ino)
	binary.BigEndian.PutUint32(data[24:28], e.mode)
	binary.BigEndian.PutUint32(data[28:32], e.uid)
	binary.BigEndian.PutUint32(data[32:36], e.gid)
	binary.BigEndian.PutUint32(data[36:40], e.file_size)
	fmt.Printf("blob oid: %x\n", e.oid)
	copy(data[40:60], e.oid)
	fmt.Println("flags: ", e.flags)
	// Writing flags to data
	binary.BigEndian.PutUint16(data[60:62], e.flags)
	// Writing path to data
	copy(data[62:], []byte(e.path))
	// Printing data path
	data = append(data, []byte(fmt.Sprintf("\x00"))...)

	for len(data)%ENTRY_BLOCK != 0 {
		data = append(data, []byte(fmt.Sprintf("\x00"))...)
	}

	// Encode the data
	return string(data)
}

func (e Entry) Parse(data []byte) Entry {
	e.ctime = int64(binary.BigEndian.Uint32(data[:4]))
	e.ctime_nsec = int64(binary.BigEndian.Uint32(data[4:8]))
	e.mtime = int64(binary.BigEndian.Uint32(data[8:12]))
	e.mtime_nsec = int64(binary.BigEndian.Uint32(data[12:16]))
	e.dev = binary.BigEndian.Uint32(data[16:20])
	e.ino = binary.BigEndian.Uint32(data[20:24])
	e.mode = binary.BigEndian.Uint32(data[24:28])
	e.uid = binary.BigEndian.Uint32(data[28:32])
	e.gid = binary.BigEndian.Uint32(data[32:36])
	e.file_size = binary.BigEndian.Uint32(data[36:40])
	e.oid = data[40:60]
	e.flags = binary.BigEndian.Uint16(data[60:62])
	e.path = strings.TrimRight(string((data[62:])), "\x00")
	return e
}

func (e *Entry) GetPath() string {
	return e.path
}
