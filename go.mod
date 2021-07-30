module github.com/splunk/jupyterhub-istio-proxy

go 1.15

require (
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/gin-gonic/gin v1.7.0
	github.com/golang/mock v1.4.4
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	istio.io/api v0.0.0-20200817160544-291eb3ba8ada
	istio.io/client-go v0.0.0-20200817160837-c5f8590ec455
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
)
