package main

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *istioClient) deleteRoute(name string) error {
	return i.NetworkingV1alpha3().VirtualServices(i.namespace).Delete(context.Background(), name, v1.DeleteOptions{})
}
