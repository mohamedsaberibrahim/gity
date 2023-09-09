package index

import (
	"reflect"
	"syscall"
	"testing"
)

func TestIndex_AddFileToIndex(t *testing.T) {
	index := &Index{}
	index.New("path/to/index")
	index.Add("alice.txt", []byte{1, 2, 3, 4}, syscall.Stat_t{})
	expected := []string{
		"alice.txt",
	}
	result, _ := index.GetSortedEntries()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Paths returned %v, expected %v", result, expected)
	}
}

func TestIndex_ReplaceFileWithDirectory(t *testing.T) {
	index := &Index{}
	index.New("path/to/index")

	index.Add("alice.txt", []byte{1, 2, 3, 4}, syscall.Stat_t{})
	index.Add("bob.txt", []byte{5, 6, 7, 8}, syscall.Stat_t{})

	index.Add("alice.txt/nested.txt", []byte{9, 10, 11, 12}, syscall.Stat_t{})

	expected := []string{
		"alice.txt/nested.txt",
		"bob.txt",
	}
	result, _ := index.GetSortedEntries()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Paths returned %v, expected %v", result, expected)
	}
}

func TestIndex_ReplaceDirectoryWithFile(t *testing.T) {
	index := &Index{}
	index.New("path/to/index")

	index.Add("alice.txt", []byte{1, 2, 3, 4}, syscall.Stat_t{})
	index.Add("nested/bob.txt", []byte{5, 6, 7, 8}, syscall.Stat_t{})

	index.Add("nested", []byte{9, 10, 11, 12}, syscall.Stat_t{})

	expected := []string{
		"alice.txt", "nested",
	}
	result, _ := index.GetSortedEntries()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Paths returned %v, expected %v", result, expected)
	}
}

func TestIndex_RecursivelyReplaceDirectoryWithFile(t *testing.T) {
	index := &Index{}
	index.New("path/to/index")

	index.Add("alice.txt", []byte{1, 2, 3, 4}, syscall.Stat_t{})
	index.Add("nested/bob.txt", []byte{5, 6, 7, 8}, syscall.Stat_t{})
	index.Add("nested/inner/claire.txt", []byte{15, 16, 7, 8}, syscall.Stat_t{})

	index.Add("nested", []byte{9, 10, 11, 12}, syscall.Stat_t{})

	expected := []string{
		"alice.txt",
		"nested",
	}
	result, _ := index.GetSortedEntries()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Paths returned %v, expected %v", result, expected)
	}
}
