# jupyterhub-istio-proxy

`jupyterhub-istio-proxy` is a scalable solution for Jupyterhub's high network traffic demands. It implements the jupyterhub proxy api to configure istio based on requests from hub.

The following requests are supported:
1. `GET /api/routes`: Gets all routes that have been configured on istio
2. `POST /api/routes/<path>`: Add the route to istio
3. `DELETE /api/routes/<path>`: Remove the route from istio

Since the proxy is stateless, it can be scaled horizontally. Multiple replicas can be used to ensure uptime during deployments and handle pod failure.

## Prerequisites

In order to use the `jupyterhub-istio-proxy` the following prerequisites need to be met.
1. The Kubernetes cluster should have istio enabled.
2. If an [istio gateway](https://istio.io/latest/docs/reference/config/networking/gateway/) is used, it should be setup to handle traffic for the FQDN where the jupyterhub instance is exposed.
3. The service account used for deploying the `jupyterhub-istio-proxy` should have ability to list, get, create and delete istio virtual services in the namespace where the deployment is done. Refer [Kubernetes RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole) for details.

### Difference to the configurable-http-proxy

Unlike the default `configurable-http-proxy` that ships with Jupyterhub, the traffic is not routed through the proxy itself. The proxy configures istio to handle all traffic to the notebook servers as well as Jupyterhub. As a result, the `proxy-public` service is not needed when using `jupyterhub-istio-proxy`. For more information see https://medium.com/@harsimran.maan/running-jupyterhub-with-istio-service-mesh-on-kubernetes-a-troubleshooting-journey-707039f36a7b

## Deployment

The proxy can be deployed to a Kubernetes namespace running Jupyterhub by applying the following config:
Change SUB_DOMAIN_HOST to a value to a hostname where jupyterhub is hosted. The ISTIO_GATEWAY value should be set to
the gateway which handles traffic for jupyterhub.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: jupyterhub
    component: proxy
  name: proxy
spec:
  replicas: 3
  selector:
    matchLabels:
      app: jupyterhub
      component: proxy
      release: RELEASE-NAME
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: jupyterhub
        component: proxy
        release: RELEASE-NAME
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: name
                  operator: In
                  values:
                  - proxy
              topologyKey: kubernetes.io/hostname
            weight: 100
      containers:
        - command:
            - /proxy/jupyterhub-istio-proxy
          env:
            - name: CONFIGPROXY_AUTH_TOKEN
              valueFrom:
                secretKeyRef:
                  key: proxy.token
                  name: hub-secret
            - name: ISTIO_GATEWAY
              value: jupyterhub-gateway
            - name: K8S_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: SUB_DOMAIN_HOST
              value: '*'
            - name: VIRTUAL_SERVICE_PREFIX
              value: jupyterhub
            - name: WAIT_FOR_WARMUP
              value: "true"
          image: splunk/jupyterhub-istio-proxy:0.0.2
          imagePullPolicy: IfNotPresent
          name: proxy
          ports:
            - containerPort: 8000
              name: proxy-api
              protocol: TCP
          resources:
            limits:
              cpu: "1"
              memory: 256M
            requests:
              cpu: 100m
              memory: 256M
          securityContext:
            allowPrivilegeEscalation: false
      securityContext:
        runAsNonRoot: true
      terminationGracePeriodSeconds: 60
---
apiVersion: v1
kind: Service
metadata:
  name: proxy-api
spec:
  ports:
    - name: http-proxy-api
      port: 8001
      protocol: TCP
      targetPort: 8000
  selector:
    component: proxy
  type: ClusterIP
---
```

Jupyterhub user pod creation flow when using `jupyterhub-istio-proxy`.
![jupyterhub-istio-proxy](http://www.plantuml.com/plantuml/png/jPD1IyGm48Nl-HN3tkjUPG-onOkB5r74ewJDYD6r2PF9Ol-zq-bkR2k227jgoVlUPDwZtIQsnFbZR-gM0y5ZGWAR8ClJH95ywwFj65OtkLaDocjkvi9RZZqZoNdb4_jGHGgVlRBwzcoZdpjkSuFK8ME2-cwdvFjbKiuOcGFLrRTr0xLpi8Rxa1bDEHP6JONGk-7WARFTGq8w-1RXHJAjpH5SpDsT79mdZfRGChfoqzBrP3sFOu66bO03D0ZYSltS94asXJgD7OejuiDGldQjroDf-ccoQxMDI0nkr7y2ixm57g6oIn7AChzqhTp_-bRlUVhMqN_hV48keWwAzAvbWtxxw2w0q7d2bcNkCS4MEoT_nHS0)

# Testing setup

https://github.com/golang/mock is used for creating mocks for testing.

```bash
mockgen --source=istio.go -destination=istio_mock_test.go -write_package_comment -package=main
```
