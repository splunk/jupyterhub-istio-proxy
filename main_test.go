package main

import "testing"

func TestValidateRequired(t *testing.T) {
	err := validateRequired("a", "val")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	err = validateRequired("a", "")
	if err == nil {
		t.Fatal("Expected error but found nil")
	}
	expected := "missing required param a"
	if expected != err.Error() {
		t.Errorf("expected %q, found %q", expected, err.Error())
	}
}
