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
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

type route struct {
	//the route specification ([host]/path/)
	RouteSpec string `json:"route_spec"`
	//the target host URL (proto://host) for this route
	Target string `json:"target"`
	// Whether the route is a JH route
	Jupyterhub bool `json:"jupyterhub,omitempty"`
	//the attached data dict for this route
	Data map[string]interface{} `json:"data,omitempty"` // Rest of the fields should go here.
}

func (r *route) splitTarget() (hostname string, port uint32) {
	u, err := url.Parse(r.Target)
	if err != nil {
		return "", port
	}

	p, _ := strconv.ParseUint(u.Port(), 10, 32)
	return u.Hostname(), uint32(p)
}

func unmarshalRoute(name string, in io.Reader) (*route, error) {
	var r route
	var m map[string]interface{}
	if err := json.NewDecoder(in).Decode(&m); err != nil {
		return nil, err
	}
	if n, ok := m["target"].(string); ok {
		r.Target = n
	}
	if n, ok := m["jupyterhub"].(bool); ok {
		r.Jupyterhub = n
	}
	r.RouteSpec = name
	delete(m, "jupyterhub")
	delete(m, "target")
	delete(m, "routespec")
	r.Data = m
	return &r, nil
}

func marshalRoute(r route) (name string, body map[string]interface{}) {
	var m = make(map[string]interface{})
	name = r.RouteSpec
	m["target"] = r.Target
	m["jupyterhub"] = r.Jupyterhub
	if r.Data != nil {
		for k, v := range r.Data {
			m[k] = v
		}
	}
	return name, m
}

func encodeRoute(r route) (string, error) {
	var b, err = json.Marshal(r)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func decodeRoute(s string) (*route, error) {
	var b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var r = route{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
