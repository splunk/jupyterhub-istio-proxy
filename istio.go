package main

import (
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/rest"
)

const (
	jupyterhubAnnotationConstant = "jupyter.splunk.io/proxy-data"
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
