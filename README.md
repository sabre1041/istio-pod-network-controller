Istio Pod Network Controller
========================

Controller to manage Istio Pod Network

This controller will initialize the pod iptables rules so that the istio proxy will inetrcept the correct connections.
This controller helps alleviate a security issue of istio. 
Without this controller, pods in the mesh must be privileged.
With this controller pods in the mesh can run with much lower privileges. 
They just have to be able to run with a defined UID which is not root.
In OpenShift this corresponds moving from the `privileged` scc to the `nonroot` scc.
The ability to run with a specific uid is required by the stio proxy.

## How this works

This controller is deployed as a daemonset
Each pod of this Daemonset takes care of the pdos deployed in the respective node.
Each pod of this daemon watches for newly created pods, if they belongs to the istio mesh, it configures the iptables of the pod so to make it join the mesh.
It then marks the pods as initialized with an annotation.

You can find instrictions on how to build this project [here](./build.html)

## Installation

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
oc apply -f applier/templates/istio-demo.yaml
oc expose svc istio-ingressgateway
oc expose svc servicegraph
oc expose svc grafana
oc expose svc prometheus
oc expose svc tracing
```

### Install istio-pod-network-controller
The istio-pod-network-controller will be installed in the istio-system namespace together with the other istio components
to install the istio-pod-network-controller run the following commands
```
oc process -f applier/templates/policies.yml NAMESPACE=istio-system | oc apply -f -
oc adm policy add-scc-to-user privileged -z istio-pod-network-controller -n istio-system
oc process -f applier/templates/istio-pod-network-controller-daemonset.yml INCLUDE_NAMESPACES=bookinfo  IMAGE=quay.io/raffaelespazzoli/istio-pod-network-controller:latest | oc apply -f - -n istio-system
oc delete cm istio-sidecar-injector -n istio-system
oc create configmap istio-sidecar-injector --from-file=config=applier/templates/istio-sidecar-injector.txt -n istio-system
```

Note that we are configuring the controller to only scan for pods in the `bookinfo` namespace.
Also note that we have modified the side car injection config map to not inject the initcontainer which is the one that requires privileged access.

## Testing with the bookinfo app

We are going to test the istio-pod-network-controller with the classic bookinfo app.

### Testing with manual sidecar injection

Run the following commands:
```
oc new-project bookinfo
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f <(istioctl kube-inject --injectConfigMapName istio-sidecar-injector -f applier/templates/bookinfo.yaml) -n bookinfo
oc expose svc productpage -n bookinfo
```

At this point the pods should deploy correclty and should participate to the istio mesh.

### Testing with automatic sidecar injection (not working in OpenShift) 

Run the following commands:
```
oc new-project bookinfo
oc label namespace bookinfo istio-injection=enabled
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f applier/templates/bookinfo.yaml -n bookinfo
```

At this point the pods should deploy correclty and should participate to the istio mesh.

