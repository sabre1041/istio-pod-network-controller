Istio Pod Network Controller
========================

Controller to manage Istio Pod Network

## Overview

This controller emulates the functionality of the [Istio init proxy](https://github.com/istio/init) to modify the _iptables_ rules so that the [Istio proxy](https://hub.docker.com/r/istio/proxyv2/) sidecar will properly intercept connections.

The primary benefit of this controller is that it helps alleviate a security issue of Istio which requires pods within the mesh to be running as privileged. Instead, privileged actions are performed by the controller instead of pods deployed by regular users. In OpenShift, this avoids the use of the `privileged` [Security Context Constraint](https://docs.openshift.com/container-platform/latest/admin_guide/manage_scc.html) and using a more restrictive policy, such as `restricted`.

## How this works

This controller is deployed as a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) that runs on each node. Each pod deployed by the DaemonSet takes on the responsibility of managing the pods that are deployed on the respective nodes the controller is deployed on.

As new pods that are to be added to the Istio mesh are created, the controller modifies iptables rules on the nodes so that the pod is able to join the mesh. Finally, the controller annotates the pod indicating that it has been successfully initialized. 

## Installation

Use the following steps to Install Istio and the controller

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
oc process -f applier/templates/policies.yml NAMESPACE=istio-system | oc apply -f -
oc adm policy add-scc-to-user privileged -z istio-pod-network-controller -n istio-system
oc process -f applier/templates/istio-pod-network-controller-daemonset.yml INCLUDE_NAMESPACES=bookinfo  IMAGE=quay.io/raffaelespazzoli/istio-pod-network-controller:latest | oc apply -f - -n istio-system
oc delete cm istio-sidecar-injector -n istio-system
oc create configmap istio-sidecar-injector --from-file=config=applier/templates/istio-sidecar-injector.txt -n istio-system
```

Take note that the controller is managing pods that are deployed in the `bookinfo` namespace as noted by he `INCLUDE_NAMESPACES` parameter of the template. This will  configure an environment in the resulting DaemonSet.

In addition, the default side car injection ConfigMap has been modified to remove the execution normally performed by the initcontainer which is by the `istioctl` tool or automatic injection process to not inject the initcontainer that typically requires privileged access.

## Testing with the bookinfo Application

To demonstrate the functionality of the `istio-pod-network-controller`, let's use he classic bookinfo application.

### Testing with manual sidecar injection

Execute the following commands:

```
oc new-project bookinfo
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f <(istioctl kube-inject -f applier/templates/bookinfo.yaml) -n bookinfo
oc expose svc productpage -n bookinfo
```

### Testing with automatic sidecar injection (not currently working in OpenShift) 

Execute the following commands:
```
oc new-project bookinfo
oc label namespace bookinfo istio-injection=enabled
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f applier/templates/bookinfo.yaml -n bookinfo
```

Regardless of the method of injection, the pods should have deployed successfully and participate in the Istio service mesh.

## Building

Instructions for building this project can be found [here](./build.md)