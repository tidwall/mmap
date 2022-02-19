package mmap

import (
	"os"
	"testing"
)

func TestMMap(t *testing.T) {
	defer os.RemoveAll("test.dat")
	err := os.WriteFile("test.dat", []byte("hello world"), 0666)
	if err != nil {
		t.Fatal(err)
	}
	data, err := Open("test.dat", false)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Fatalf("expected '%s' got '%s'", "hello world", string(data))
	}
	if err := Close(data); err != nil {
		t.Fatal(err)
	}
	data, err = Open("test.dat", true)
	if err != nil {
		t.Fatal(err)
	}
	data[0] = 'j'
	copy(data[6:], "earth")
	if string(data) != "jello earth" {
		t.Fatalf("expected '%s' got '%s'", "jello world", string(data))
	}
	if err := Close(data); err != nil {
		t.Fatal(err)
	}
	data, err = Open("test.dat", false)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "jello earth" {
		t.Fatalf("expected '%s' got '%s'", "jello earth", string(data))
	}
	if err := Close(data); err != nil {
		t.Fatal(err)
	}
}
