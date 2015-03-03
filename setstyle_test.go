package main

import (
	"testing"
)

func TestSetStyleTrue(t *testing.T) {
	actual := setStyle("vim")
	expected := "vim"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestSetStyleFalse(t *testing.T) {
	actual := setStyle("unknown")
	expected := "default"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
