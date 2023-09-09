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
