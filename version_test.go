package main

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	actual := versionInfo()
	expected := "Version: dev Commit: HEAD GoVersion: Unknown"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}
