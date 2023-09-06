package database

import (
	"fmt"
	"sort"

	"github.com/mohamedsaberibrahim/gity/index"
)

type Tree struct {
	entries map[string]ObjectInterface
	oid     []byte
}

func (t *Tree) New() {
	t.entries = map[string]ObjectInterface{}
}

func (t *Tree) ToString() string {
	var entries []byte

	for name, entry := range t.entries {
		entries = append(entries, []byte(fmt.Sprintf("%s %s%x", entry.GetMode(), name, 0x00))...)

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

func (t *Tree) GetType() string {
	return TREE_TYPE
}

func (t *Tree) GetMode() string {
	return DIRECTORY_MODE
}

func (t *Tree) AddEntry(parents []string, e Entry) {
	if len(parents) == 0 {
		t.entries[GetBaseName(e.ToString())] = ObjectInterface(&e)
	} else {
		firstParent := parents[0]
		value, ok := t.entries[firstParent]
		var tree Tree
		if ok {
			tree = *(value.(*Tree))
		} else {
			tree.New()
		}
		t.entries[GetBaseName(firstParent)] = ObjectInterface(&tree)
		t.AddEntry(parents[1:], e)
	}
}

func (t *Tree) Traverse(callback func(object ObjectInterface) error) {
	for _, value := range t.entries {
		if value.GetType() == TREE_TYPE {
			value.(*Tree).Traverse(callback)
		}
	}

	callback(ObjectInterface(t))
}

func (t Tree) Build(entries []index.Entry) Tree {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].ToString() < entries[i].ToString()
	})

	root := Tree{}
	root.New()

	for _, indexEntry := range entries {
		parents := indexEntry.GetParentDirectories("")
		databaseEntry := Entry{}
		databaseEntry.New(indexEntry.ToString(), indexEntry.GetOid(), indexEntry.GetStat())
		root.AddEntry(parents, databaseEntry)
	}

	return root
}
