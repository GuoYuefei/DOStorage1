package es

import (
	"fmt"
	"testing"
)

func TestSearchLatestVersion(t *testing.T) {
	meta, err := SearchLatestVersion("somefd")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(meta)
}

func TestGetMetadata(t *testing.T) {
	metadata, err := GetMetadata("abc", 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(metadata)
}
