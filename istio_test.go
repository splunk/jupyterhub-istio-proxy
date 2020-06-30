package main

import (
	"testing"
)

func TestVirtualServiceAnnotationNameWithPrefix(t *testing.T) {
	actual := virtualServiceAnnotationNameWithPrefix()
	expected := "jupyter.splunk.io/proxy-data"
	if expected != actual {
		t.Errorf("expected %q, found %q. Changing this expectation would require manual recycling of exisitng VS in a deployment", expected, actual)
	}
}
