FROM gcr.io/distroless/static-debian10:nonroot

ONBUILD RUN echo -e "\
========================================================================\n\
=                                                                      =\n\
= ⚠️ DEPRECATION NOTICE                                                =\n\
= As of June 24th, 2025: (Latest Release 0.3.0).                       =\n\
=                                                                      =\n\
= Timeline:                                                            =\n\
= ~ 60 days for GitHub Code Archive ->                                 =\n\
=   https://github.com/splunk/jupyterhub-istio-proxy                   =\n\
= ~ 30 days for DockerHub Image Removal ->                             =\n\
=   https://hub.docker.com/repository/docker/splunk/                   =\n\
=   jupyterhub-istio-proxy/general                                     =\n\
=                                                                      =\n\
= Maintenance:                                                         =\n\
= Anyone actively using this code please Fork it.                      =\n\
= Anyone interested in maintaining the Repository,                     =\n\
=   raise a Pull Request for CODEOWNERS.                               =\n\
= We will then proceed to review the request internally.               =\n\
=                                                                      =\n\
========================================================================"

COPY jupyterhub-istio-proxy /proxy/jupyterhub-istio-proxy
ENTRYPOINT ["/proxy/jupyterhub-istio-proxy"]
CMD [ "/proxy/jupyterhub-istio-proxy" ]
