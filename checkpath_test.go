package main

import (
	"testing"
)

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func TestCheckPath(t *testing.T) {
	try, _ := ioutil.TempDir("/tmp", "tmp")
	os.Chdir(try)
	_, err := os.OpenFile("my_command", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		fmt.Println(err)
	}
	dir, _ := filepath.Abs(try)
	os.Setenv("PATH", os.Getenv("PATH")+":"+dir)

	actual := path.Base(checkPath("my_command"))
	expected := "my_command"

	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
