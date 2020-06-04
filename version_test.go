package main

import (
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	actual := getVersionInfo()
	expected := "Version: dev Commit: HEAD GoVersion: Unknown"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}
