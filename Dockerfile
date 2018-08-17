#FROM registry.access.redhat.com/rhel7/rhel:7.5
FROM centos:7

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
    INSTALL_PKGS="golang iptables iproute git runc" && \
    yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    mkdir -p ${GOPATH}/{bin,src} && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    cd /opt/app-root/go/src/github.com/sabre1041/istio-pod-network-controller && \
    cp bin/istio-iptables.sh /usr/local/bin/ && \
    dep ensure -vendor-only &&\
    go build -o bin/istio-pod-network-controller -v cmd/istio-pod-network-controller/main.go && \
    mv bin/istio-pod-network-controller /usr/local/bin && \
    rm -rf ${GOPATH} && \
    REMOVE_PKGS="golang git" && \
    yum remove -y $REMOVE_PKGS && \
    yum clean all && \
    rm -rf /var/cache/yum && \
    VERSION="v1.11.1" && \
    curl -L -o /root/crictl-$VERSION-linux-amd64.tar.gz https://github.com/kubernetes-incubator/cri-tools/releases/download/$VERSION/crictl-$VERSION-linux-amd64.tar.gz &&\
    tar zxvf /root/crictl-$VERSION-linux-amd64.tar.gz -C /usr/local/bin && \
    rm -f crictl-$VERSION-linux-amd64.tar.gz && \
    curl -L -o /usr/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
    chmod +x /usr/bin/jq   


ENTRYPOINT ["/usr/local/bin/istio-pod-network-controller"]