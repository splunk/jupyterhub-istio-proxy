# jupyterhub-istio-proxy

`jupyterhub-istio-proxy` is a jupyterhub proxy api implementation that is responsible for configuring istio based on requests from hub.

The following requests are supported:
1. `GET /api/routes`: Gets all routes that have configured on istio
2. `POST /api/routes/<path>`: Add the route to istio
3. `DELETE /api/routes/<path>`: Remove the route from istio

Since the proxy is stateless, it can be scaled horizontally. Multiple replicas can be used to ensure uptime during deployments and handle pod failure.

# Testing setup

https://github.com/golang/mock is used for creating mocks for testing.

```bash
mockgen --source=istio.go -destination=istio_mock_test.go -write_package_comment -package=main
```
