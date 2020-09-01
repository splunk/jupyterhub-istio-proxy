# jupyterhub-istio-proxy

`jupyterhub-istio-proxy` is a jupyterhub proxy api implementation that is responsible for configuring istio based on requests from hub.

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

Unlike the default `configurable-http-proxy` that ships with istio, the traffic is not routed through the proxy itself. The proxy configures istio to handle all traffic to the notebook servers as well as Jupyterhub. As a result, the `proxy-public` service is not needed when using `jupyterhub-istio-proxy`

## Deployment

The proxy can be deployed to a Kubernetes namespace running Jupyterhub by applying the following config:
Change SUB_DOMAIN_HOST to a value to a hostname where jupyterhub is hosted. The ISTIO_GATEWAY value should be set to
the gateway which handles traffic for jupyterhub.

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: jupyterhub
    component: proxy
  name: proxy
spec:
  replicas: 3
  revisionHistoryLimit: 1
  selector:
    matchLabels:
      app: jupyterhub
      component: proxy
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: jupyterhub
        component: proxy
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
          value: istio-gateway/jupyterhub-gateway
        - name: K8S_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: SUB_DOMAIN_HOST
          value: example.com
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
      nodeSelector: {}
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
    port: 80
    protocol: TCP
    targetPort: 8000
  selector:
    component: proxy
  type: ClusterIP

```
Jupyterhub user pod creation flow when using `jupyterhub-istio-proxy`.
![jupyterhub-istio-proxy](http://www.plantuml.com/plantuml/png/RP4nJmCn38Nt_8gdxZ2Z6r1FB5Gi5GWnejp5Y8XoYHodel-UE8IhGfbEv_UzPyVU9h4i-VDWnba2upaHmoRayZMnI7xsqIw2pNsUDggyvwaNzXo-JXZtkof7_NkrqVVGGvw85n9AA_bnaofdj1UkRZLvm9FEWr8v4fjIKQ6H0-wOh11YS3_IFvew_KocrTVSui5S4TTYpwCE69t3OyX2PYrsZLNtINp0qNCMUxZeq-VjVbd9PtRjty0BrulrWhV0O1q54Z0HYiUTUei-g35_dgIhB3kisuEUjcZNO7AUQHPDLA8kS4kBAtV_1W00)

# Testing setup

https://github.com/golang/mock is used for creating mocks for testing.

```bash
mockgen --source=istio.go -destination=istio_mock_test.go -write_package_comment -package=main
```
