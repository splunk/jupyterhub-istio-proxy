package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"istio.io/api/networking/v1alpha3"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jupyterProxyIDHeaderKey string = "jupyter-proxy-id"
	maxRetries              uint64 = 6
)

func (i *istioClient) createVirtualService(r route) error {
	log.Println("creating route: ", r)
	var annotations, err = getAnnotationForRoute(r)
	if err != nil {
		return err
	}
	destinationHost, destinationPort := r.splitTarget()
	destinationHost = fmt.Sprintf("%s.%s.svc.cluster.local", destinationHost, i.namespace)
	vsName := getVirtualServiceNameWithPrefix(r.RouteSpec)
	vs := getVirtualService(vsName, i.gateway, i.host, destinationHost, destinationPort, r.RouteSpec, annotations)

	_, err = i.NetworkingV1alpha3().VirtualServices(i.namespace).Create(context.Background(), vs, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("virtual service creation skipped as %s already exists\n", vsName)
			return nil
		}
		log.Println("virtual service creation failed", err)
		return err
	}
	log.Println("virtual service created")
	if i.waitForWarmup {
		log.Println("waiting for warmup")
		err = warmup(vsName, getWarmupURL(i.host, r.RouteSpec))
	}
	if err != nil {
		log.Printf("warming up the servers did not return after %d tries. Continuing despite: %v", maxRetries+1, err)
		return nil
	}
	return nil
}

func getVirtualServiceNameWithPrefix(name string) string {
	sum := sha256.Sum256([]byte(name))
	return fmt.Sprintf("%s-%x", getVsNamePrefix(), sum)
}

func getVsNamePrefix() string {
	if vsNamePrefix != "" {
		return vsNamePrefix
	}
	return vsNamePrefixDefault
}

func getWarmupURL(host string, p string) string {
	return fmt.Sprintf("https://%s/", path.Join(host, p))
}

func getAnnotationForRoute(r route) (map[string]string, error) {
	var e, err = encodeRoute(r)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]string)
	m[jupyterhubAnnotationConstant] = e
	return m, nil
}

func warmup(name string, url string) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	fetchURL := func() error {
		log.Println("GETing url to check if it is up", url)
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("error GETing ", url, err)
			return err
		}
		defer resp.Body.Close()
		if name == resp.Header.Get(jupyterProxyIDHeaderKey) {
			log.Println("virtual service is warmed up")
			return nil
		}
		return fmt.Errorf("desired header `%s` not found", name)
	}
	bf := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
	return backoff.Retry(fetchURL, bf)
}

func getVirtualService(name string, gateway string, host string, destinationHost string, destinationPort uint32, route string, annotations map[string]string) *networkingv1alpha3.VirtualService {
	return &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
		},
		Spec: v1alpha3.VirtualService{
			Hosts:    []string{host},
			Gateways: []string{gateway},
			Http: []*v1alpha3.HTTPRoute{
				{
					Match: []*v1alpha3.HTTPMatchRequest{
						{
							Uri: &v1alpha3.StringMatch{
								MatchType: &v1alpha3.StringMatch_Prefix{
									Prefix: route,
								},
							},
						},
					},
					Route: []*v1alpha3.HTTPRouteDestination{
						{
							Destination: &v1alpha3.Destination{
								Host: destinationHost,
								Port: &v1alpha3.PortSelector{
									Number: destinationPort,
								},
							},
						},
					},
					Headers: &v1alpha3.Headers{
						Response: &v1alpha3.Headers_HeaderOperations{
							Set: map[string]string{
								jupyterProxyIDHeaderKey: name,
							},
						},
					},
				},
			},
		},
	}
}
