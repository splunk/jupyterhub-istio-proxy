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
	"fmt"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/rest"
)

type istioer interface {
	createVirtualService(route) error
	listRegisteredRoutes() (map[string]interface{}, error)
	deleteRoute(string) error
}

type istioClient struct {
	*versionedclient.Clientset
	gateway       string
	host          string
	namespace     string
	waitForWarmup bool
}

func newIstioClient(namespace string, gateway string, host string, waitForWarmup bool) (*istioClient, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	ic, err := versionedclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &istioClient{Clientset: ic, namespace: namespace, gateway: gateway, host: host, waitForWarmup: waitForWarmup}, nil
}
func virtualServiceAnnotationNameWithPrefix() string {
	return fmt.Sprintf("%s.splunk.io/proxy-data", virtualServicePrefix())
}
