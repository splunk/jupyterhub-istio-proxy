/*
Copyright 2020 Splunk Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"testing"
)

func TestVirtualServiceNameWithPrefix(t *testing.T) {
	var name = "/"
	actual := virtualServiceNameWithPrefix(name)
	expected := "jupyter-8a5edab282632443219e051e4ade2d1d5bbc671c781051bf1437897cbdfea0f1"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}

func TestWarmupURL(t *testing.T) {
	actual := warmupURL("example.com", "/my-path")
	expected := "https://example.com/my-path/"
	if expected != actual {
		t.Errorf("expected %q, found %q", expected, actual)
	}
}

func TestVirtualService(t *testing.T) {
	name := "my-vs"
	gateway := "hsm/gateway"
	host := "example.com"
	service := "mylocalservice"
	port := uint32(80)
	path := "/my-path"
	annotations := map[string]string{"test/annotation": "test"}
	vs := virtualService(name, gateway, host, service, port, path, annotations)

	if name != vs.Name {
		t.Errorf("expected %q, found %q", name, vs.Name)
	}
	if annotations["test/annotation"] != vs.Annotations["test/annotation"] {
		t.Errorf("expected %q, found %q", annotations, vs.Annotations)
	}

}
