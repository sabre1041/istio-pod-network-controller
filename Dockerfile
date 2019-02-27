#FROM registry.access.redhat.com/rhel7/rhel:7.5
FROM centos:7

LABEL io.k8s.description="Image to Manage Istio's Pod Network" \
      io.k8s.display-name="Istio Pod Network Controller" \
      io.openshift.tags="go"

RUN yum repolist > /dev/null && \
    yum-config-manager --enable rhel-7-server-optional-rpms --enable rhel-7-server-extras-rpms && \
    yum clean all && \
    INSTALL_PKGS="iptables iproute runc" && \
    yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all && \
    rm -rf /var/cache/yum && \
    VERSION="v1.11.1" && \
    curl -L -o /root/crictl-$VERSION-linux-amd64.tar.gz https://github.com/kubernetes-incubator/cri-tools/releases/download/$VERSION/crictl-$VERSION-linux-amd64.tar.gz &&\
    tar zxvf /root/crictl-$VERSION-linux-amd64.tar.gz -C /usr/local/bin && \
    rm -f crictl-$VERSION-linux-amd64.tar.gz && \
    curl -L -o /usr/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
    chmod +x /usr/bin/jq   

ADD ./bin/istio-pod-network-controller /usr/local/bin/istio-pod-network-controller
ADD ./bin/istio-iptables.sh /usr/local/bin/istio-iptables.sh

ENTRYPOINT ["/usr/local/bin/istio-pod-network-controller"]
