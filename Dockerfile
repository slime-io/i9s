FROM ubuntu:bionic
ARG KUBECTL_VERSION="v1.21.2"
ARG FX_ARCH="amd64"
COPY ./execs/i9s /bin/i9s
COPY ./bin/istioctl /bin/istioctl
RUN apt-get update && apt-get install -y jq less vim curl \
    && chmod +x /bin/istioctl \
    && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    && curl -L https://github.com/antonmedv/fx/releases/download/24.0.0/fx_linux_${FX_ARCH} -o /usr/local/bin/fx \
    && chmod +x /usr/local/bin/fx
ENTRYPOINT [ "/bin/i9s" ]