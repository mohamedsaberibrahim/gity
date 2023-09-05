package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	CHECKSUM_SIZE = 20
)

type Checksum struct {
	file_data []byte
	file      *os.File
	offset    int64
}

func (c *Checksum) New(file *os.File) {
	c.file = file
	c.file_data = []byte{}
}

func (c *Checksum) Read(size int) ([]byte, error) {
	fmt.Println("Trying to read bytes: ", size)
	fmt.Println("offset: ", c.offset)
	// if size == 8 {
	// 	size = 5
	// }
	data := make([]byte, size)
	n_bytes, err := c.file.Read(data)
	c.file_data = append(c.file_data, data...)
	c.offset += int64(n_bytes)

	c.file.Seek(c.offset, 0)
	fmt.Println("n_bytes: ", n_bytes)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Checksum) VerifyChecksum() bool {
	c.file.Seek(-CHECKSUM_SIZE, 2)
	var oid [20]byte
	c.file.Read(oid[0:20])
	fmt.Printf("oid: %x\n", oid)

	var encoded_data bytes.Buffer
	if err := binary.Write(&encoded_data, binary.BigEndian, c.file_data); err != nil {
		fmt.Printf("Error: failed to write data - %v\n", err)
		return false
	}
	data_objectId := sha1.Sum(encoded_data.Bytes())
	fmt.Printf("data_objectId: %x\n", data_objectId)
	if data_objectId != oid {
		fmt.Println("Checksums do not match")
		return false
	}
	return true
}

func (c *Checksum) Write(data []byte) {
	fmt.Println("Before Writing data: ", len(c.file_data))
	c.file.Write(data)
	c.file_data = append(c.file_data, data...)
	fmt.Println("After Writing data: ", len(c.file_data))
}

func (c *Checksum) WriteChecksum() {
	fmt.Println("Final data size: ", len(c.file_data))
	var encoded_data bytes.Buffer
	if err := binary.Write(&encoded_data, binary.BigEndian, c.file_data); err != nil {
		fmt.Printf("Error: failed to write data - %v\n", err)
	}
	data_objectId := sha1.Sum(encoded_data.Bytes())
	fmt.Printf("data_objectId: %x\n", data_objectId)
	c.file.Write(data_objectId[:])
}
