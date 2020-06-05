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
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var sharedAPIToken string
var subDomainHost string
var gateway string
var namespace string
var waitForWarmup bool
var vsNamePrefix string

const (
	gatewayEnvKey        = "ISTIO_GATEWAY"
	proxyAuthTokenEnvKey = "CONFIGPROXY_AUTH_TOKEN"
	subDomainHostEnvKey  = "SUB_DOMAIN_HOST"
	namespaceKey         = "K8S_NAMESPACE"
	waitForWarmupKey     = "WAIT_FOR_WARMUP"
	vsNamePrefixKey      = "VS_NAME_PREFIX"
	vsNamePrefixDefault  = "jupyter"
)

func main() {
	log.Println(getVersionInfo())
	gin.SetMode(gin.ReleaseMode)

	sharedAPIToken = os.Getenv(proxyAuthTokenEnvKey)
	err := validateRequired(proxyAuthTokenEnvKey, sharedAPIToken)
	if err != nil {
		log.Fatalln(err)
	}

	gateway = os.Getenv(gatewayEnvKey)
	err = validateRequired(gatewayEnvKey, gateway)
	if err != nil {
		log.Fatalln(err)
	}
	namespace = os.Getenv(namespaceKey)
	err = validateRequired(namespaceKey, namespace)
	if err != nil {
		log.Fatalln(err)
	}
	subDomainHost = os.Getenv(subDomainHostEnvKey)
	err = validateRequired(subDomainHostEnvKey, subDomainHost)
	if err != nil {
		log.Fatalln(err)
	}
	waitForWarmup = os.Getenv(waitForWarmupKey) != "false"
	var ok bool
	if vsNamePrefix, ok = os.LookupEnv(vsNamePrefixKey); !ok {
		vsNamePrefix = vsNamePrefixDefault
	}
	var ic istioer
	ic, err = newIstioClient(namespace, gateway, subDomainHost, waitForWarmup)
	if err != nil {
		log.Fatalf("failed to create istio client: %s\n", err)
	}
	r := gin.Default()
	registerRoutes(r, ic, sharedAPIToken)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// Handle signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
