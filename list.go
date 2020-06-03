package main

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *istioClient) listRegisteredRoutes() (map[string]interface{}, error) {
	vsList, err := i.NetworkingV1alpha3().VirtualServices(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var routes = make(map[string]interface{})
	for i := range vsList.Items {
		vs := vsList.Items[i]
		if a, ok := vs.Annotations[jupyterhubAnnotationConstant]; ok {
			rd, err := decodeRoute(a)
			if err != nil {
				log.Printf("error decoding annotation but continuing: %s\n", vs.Name)
				continue
			}
			name, mp := marshalRoute(*rd)
			log.Printf("Added %q to list\n", name)
			routes[name] = mp

		} else {
			log.Printf("skipping vs as it does not have the required annotation: %s\n", vs.Name)
		}
	}

	return routes, nil
}
