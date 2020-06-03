FROM golang:1.14 as builder
COPY . /proxy
WORKDIR /proxy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . && chmod +x jupyterhub-istio-proxy

FROM gcr.io/distroless/static-debian10:nonroot
COPY --from=builder /proxy/jupyterhub-istio-proxy /proxy/jupyterhub-istio-proxy
ENTRYPOINT ["/proxy/jupyterhub-istio-proxy"]
CMD [ "/proxy/jupyterhub-istio-proxy" ]
