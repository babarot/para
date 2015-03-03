package main

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	actual := stringInSlice("a", []string{"a", "b", "c"})
	expected := true
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
