package database

import (
	"fmt"
	"sort"
)

const (
	entryFormat = "Z*H40"
)

type Tree struct {
	entries []Entry
	oid     []byte
}

func (t *Tree) New(entries []Entry) {
	t.entries = entries
}

func (t *Tree) ToString() string {
	sort.SliceStable(t.entries, func(i, j int) bool {
		return t.entries[i].GetName() < t.entries[i].GetName()
	})

	var entries []byte

	for _, entry := range t.entries {
		entries = append(entries, []byte(fmt.Sprintf("%s %s%x", entry.GetMode(), entry.GetName(), 0x00))...)

		entries = append(entries, entry.GetOid()...)
	}

	return string(entries)
}

func (t *Tree) GetOid() []byte {
	return t.oid
}

func (t *Tree) SetOid(oid []byte) {
	t.oid = oid
}

func (b *Tree) GetType() string {
	return "tree"
}
