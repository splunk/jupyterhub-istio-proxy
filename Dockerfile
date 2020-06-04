FROM gcr.io/distroless/static-debian10:nonroot
COPY jupyterhub-istio-proxy /proxy/jupyterhub-istio-proxy
ENTRYPOINT ["/proxy/jupyterhub-istio-proxy"]
CMD [ "/proxy/jupyterhub-istio-proxy" ]
