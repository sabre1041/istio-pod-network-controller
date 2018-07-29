#FROM registry.access.redhat.com/rhel7/rhel:7.5
FROM centos:7

MAINTAINER Andrew Block <ablock@redhat.com>

LABEL io.k8s.description="Image to Manage Istio's Pod Network" \
      io.k8s.display-name="Istio Pod Network Controller" \
      io.openshift.tags="go"

ENV GOPATH=/opt/app-root/go \
    GOBIN=/opt/app-root/go/bin \
    PATH=$PATH:/opt/app-root/go/bin

ADD . /opt/app-root/go/src/github.com/sabre1041/istio-pod-network-controller

RUN yum repolist > /dev/null && \
#    yum-config-manager --enable rhel-7-server-optional-rpms --enable rhel-7-server-extras-rpms && \
    yum clean all && \
    INSTALL_PKGS="golang iptables iproute git" && \
    yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all && \
    mkdir -p ${GOPATH}/{bin,src} && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    cd /opt/app-root/go/src/github.com/sabre1041/istio-pod-network-controller && \
    cp bin/istio-iptables.sh /usr/local/bin/ && \
    dep ensure -vendor-only &&\
    go build -o bin/istio-pod-network-controller -v main.go && \
    mv bin/istio-pod-network-controller /usr/local/bin && \
    rm -rf ${GOPATH}

ENTRYPOINT ["/usr/local/bin/istio-pod-network-controller"]