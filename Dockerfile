FROM registry.access.redhat.com/rhel7/rhel:7.5

MAINTAINER Andrew Block <ablock@redhat.com>

LABEL io.k8s.description="Image to Manage Istio's Pod Network" \
      io.k8s.display-name="Istio Pod Network Controller" \
      io.openshift.tags="go"

ENV GOPATH=/opt/app-root/go \
    GOBIN=/opt/app-root/go/bin \
    PATH=$PATH:/opt/app-root/go/bin

ADD . /opt/app-root/go/src/github.com/sabre1041/istio-pod-network-controller

RUN yum repolist > /dev/null && \
    yum-config-manager --enable rhel-7-server-optional-rpms --enable rhel-7-server-extras-rpms && \
    yum clean all && \
    INSTALL_PKGS="golang iptables iproute docker-client grep gawk" && \
    yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all && \
    mkdir -p ${GOPATH}/{bin,src} && \
    cd /opt/app-root/go/src/github.com/sabre1041/istio-pod-network-controller && \
    cp bin/istio-iptables.sh /usr/local/bin/ && \
    go build -v -o bin/main cmd/istio-pod-network-controller/main.go && \
    mv bin/main ${GOBIN}

CMD ["/opt/app-root/go/bin/main"]