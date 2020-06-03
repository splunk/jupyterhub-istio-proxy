package main

import (
	"testing"
)

func TestGetVirtualServiceNameWithPrefix(t *testing.T) {
	var name = "/"
	actual := getVirtualServiceNameWithPrefix(name)
	expected := "jupyter-8a5edab282632443219e051e4ade2d1d5bbc671c781051bf1437897cbdfea0f1"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}

func TestGetWarmupURL(t *testing.T) {
	actual := getWarmupURL("example.com", "/my-path")
	expected := "https://example.com/my-path/"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}

func TestGetVirtualService(t *testing.T) {
	name := "my-vs"
	gateway := "hsm/gateway"
	host := "example.com"
	service := "mylocalservice"
	port := uint32(80)
	path := "/my-path"
	annotations := map[string]string{"test/annotation": "test"}
	vs := getVirtualService(name, gateway, host, service, port, path, annotations)

	if name != vs.Name {
		t.Errorf("expected %q, found %q", name, vs.Name)
	}
	if annotations["test/annotation"] != vs.Annotations["test/annotation"] {
		t.Errorf("expected %q, found %q", annotations, vs.Annotations)
	}

}
