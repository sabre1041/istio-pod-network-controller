Istio Pod Network Controller
========================

Controller to manage Istio Pod Network

## Overview

This controller emulates the functionality of the [Istio init proxy](https://github.com/istio/init) to modify the _iptables_ rules so that the [Istio proxy](https://hub.docker.com/r/istio/proxyv2/) sidecar will properly intercept connections.

The primary benefit of this controller is that it helps alleviate a security issue of Istio which requires pods within the mesh to be running as privileged. Instead, privileged actions are performed by the controller instead of pods deployed by regular users. In OpenShift, this avoids the use of the `privileged` [Security Context Constraint](https://docs.openshift.com/container-platform/latest/admin_guide/manage_scc.html) and using a more restrictive policy, such as `nonroot`.

## How this works

This controller is deployed as a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) that runs on each node. Each pod deployed by the DaemonSet takes on the responsibility of managing the pods that are deployed on the respective nodes the controller is deployed on.

As new pods that are to be added to the Istio mesh are created, the controller modifies iptables rules on the nodes so that the pod is able to join the mesh. Finally, the controller annotates the pod indicating that it has been successfully initialized. 

Pod will be initialized if the pod's namespace is annotated with `istio-pod-network-controller/initialize: true` or if the pod itself is annotated with `istio-pod-network-controller/initialize: true`. The logic works the same as for the `istio-injection: enabled` label.

## Installation on Kubernetes

### Starting Kubernetes 

If you don't have a kubernetes cluster available run this command to start a minikube instance large enough to host istio:
```
minikube start --memory=8192 --cpus=2 --kubernetes-version=v1.10.0 \
    --extra-config=controller-manager.cluster-signing-cert-file="/var/lib/localkube/certs/ca.crt" \
    --extra-config=controller-manager.cluster-signing-key-file="/var/lib/localkube/certs/ca.key"
```
If you want to run minikube with the crio container runtime run the following:
```
minikube start --memory=8192 --cpus=2 --kubernetes-version=v1.10.0 \
    --extra-config=controller-manager.cluster-signing-cert-file="/var/lib/localkube/certs/ca.crt" \
    --extra-config=controller-manager.cluster-signing-key-file="/var/lib/localkube/certs/ca.key" \
    --network-plugin=cni \
    --container-runtime=cri-o \
    --bootstrapper=kubeadm
```

### Install Istio

Run the following to install `Istio`
```
kubectl create namespace istio-system
kubectl apply -f applier/templates/istio-demo.yaml -n istio-system
```

### Install istio-pod-network-controller

Run the following to install `istio-pod-network-controller`
```
helm template -n istio-pod-network-controller ./chart/istio-pod-network-controller | kubectl apply -f -
```

if you are using with crio, run the following
```
helm template -n istio-pod-network-controller --ser containerRuntime=crio ./chart/istio-pod-network-controller | kubectl apply -f -
```

### Testing with automatic sidecar injection 

Execute the following commands:
```
kubectl create namespace bookinfo
kubectl label namespace bookinfo istio-injection=enabled
kubectl annotate namespace bookinfo istio-pod-network-controller/initialize=true
kubectl apply -f applier/templates/bookinfo.yaml -n bookinfo
```

## Installation on OpenShift

### Starting OpenShift

If you don't have an OpenShift cluster available run this command to start a minikube instance large enough to host istio:
```
minishift start --ocp-tag=v3.9.40 --vm-driver=kvm \
    --cpus=2 --memory=8192 --skip-registration

```

### Install istio

```
oc adm new-project istio-system --node-selector=""
oc adm policy add-scc-to-user anyuid -z istio-ingress-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z default -n istio-system
oc adm policy add-scc-to-user anyuid -z prometheus -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-egressgateway-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-citadel-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-ingressgateway-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-cleanup-old-ca-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-mixer-post-install-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-mixer-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-pilot-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-sidecar-injector-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-galley-service-account -n istio-system
oc apply -f applier/templates/istio-demo.yaml -n istio-system
oc expose svc istio-ingressgateway -n istio-system
oc expose svc servicegraph -n istio-system
oc expose svc grafana -n istio-system
oc expose svc prometheus -n istio-system
oc expose svc tracing -n istio-system
```

### Install istio-pod-network-controller

The _istio-pod-network-controller_ is to be installed in the `istio-system` namespace along with with the other istio components

To install the _istio-pod-network-controller_, execute the following commands:

```
helm template -n istio-pod-network-controller ./chart/istio-pod-network-controller | oc apply -f -
```

## Testing with the bookinfo Application

To demonstrate the functionality of the `istio-pod-network-controller`, let's use he classic bookinfo application.

### Testing with manual sidecar injection

Execute the following commands:

```
oc new-project bookinfo
oc annotate namespace bookinfo istio-pod-network-controller/initialize=true
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f <(istioctl kube-inject -f applier/templates/bookinfo.yaml) -n bookinfo
oc expose svc productpage -n bookinfo
```

## Building

Instructions for building this project can be found [here](./build.md)