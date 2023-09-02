package database

import (
	"fmt"
	"strings"
)

type Commit struct {
	oid      []byte
	parent   []byte
	tree_oid []byte
	author   Author
	message  string
}

func (c *Commit) New(parent []byte, tree_oid []byte, author Author, message string) {
	c.parent = parent
	c.tree_oid = tree_oid
	c.author = author
	c.message = message
}

func (c *Commit) ToString() string {
	var parent string
	if parent = ""; c.parent != nil {
		parent = fmt.Sprintf("parent %s\n", string(c.parent))
	}

	tree := fmt.Sprintf("tree %x\n", string(c.tree_oid))
	author := fmt.Sprintf("author %s\n", c.author.ToString())
	committer := fmt.Sprintf("committer %s\n", c.author.ToString())

	return strings.Join([]string{tree, parent, author, committer, "\n", c.message, "\n"}, "")
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

func (c *Commit) GetMode() string {
	return "commit"
}
