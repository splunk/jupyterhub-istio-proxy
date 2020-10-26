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

package proxy

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *IstioClient) listRegisteredRoutes() (map[string]interface{}, error) {
	vsList, err := i.NetworkingV1alpha3().VirtualServices(i.namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var routes = make(map[string]interface{})
	for _, vs := range vsList.Items {
		if a, ok := vs.Annotations[i.virtualServiceAnnotationNameWithPrefix()]; ok {
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
