package database

import (
	"fmt"
)

type Commit struct {
	oid      []byte
	tree_oid []byte
	author   Author
	message  string
}

func (c *Commit) New(tree_oid []byte, author Author, message string) {
	c.tree_oid = tree_oid
	c.author = author
	c.message = message
}

func (c *Commit) ToString() string {
	return fmt.Sprintf("tree %x\nauthor %s\ncommitter %s\n\n%s", string(c.tree_oid), c.author.ToString(), c.author.ToString(), c.message)
}

func (c *Commit) GetOid() []byte {
	return c.oid
}

func (c *Commit) SetOid(oid []byte) {
	c.oid = oid
}

func (c *Commit) GetType() string {
	return "commit"
}

func (c *Commit) GetMessage() string {
	return c.message
}
