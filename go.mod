module github.com/splunk/jupyterhub-istio-proxy

go 1.14

require (
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/mock v1.2.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	istio.io/api v0.0.0-20200512011036-83e5c7ad8375
	istio.io/client-go v0.0.0-20200512012238-d2e342e465b2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
)
