Istio Pod Network Controller
========================

Controller to manage Istio Pod Network

This controller will initialize the pod iptables rules so that the istio proxy will inetrcept the correct connections.
This controller helps alleviate a security issue of istio. 
Without this controller, pods in the mesh must be privileged.
With this controller pods in the mesh can run with much lower privileges. 
They just have to be able to run with a defined UID which is not root.
In OpenShift this corresponds to the `nonroot` scc.
The ability to run with a specific uid is required by the stio proxy.

## How this works

This controller is deployed as a daemonset
Each pod of this Daemonset takes care of the pdos deployed in the respective node.
Each pod os this daemon set watch for newly created pods, if they belongs to the istio mesh, it configures the iptables of the pod so to make it join the mesh.
The picture below describe the process.
 
TODO add picture

## Prerequisites

The following are required for the project

* OpenShift Command Line Tool
* OpenShift Account with _cluster-admin_ privileges

## Installation

Install istio (TODO Link)
The istio-pod-network-controller can run in any namespace, here we assume it's going to run in the istio-system namespace together with the other istio components
to install the istio-pod-network-controller run the following commands
```
oc process applier/policies.yml NAMESPACE=istio-system | oc apply -f -
oc adm policy add-scc-to-user privileged -z istio-pod-network-controller -n istio-system
oc process applier/istio-pod-network-controller-daemonset.yml INCLUDE_NAMESPACES=bookinfo NAMESPACE=istio-system IMAGE=quay.io/raffaelespazzoli/istio-pod-network-controller:latest | oc apply -f -
oc delete cm istio-sidecar-injector -n istio-system
oc create configmap istio-sidecar-injector --from-file=config=applier/templates/istio-sidecar-injector.txt
```

Note that we are configuring the controller to only scan for pods in the `bookinfo` namespace.
Also note that we have modified the side car injection config map to not inject the initcontainer which is the one that requires privileged access

## Testing with the bookinfo app

We are going to test the istio-pod-network-controller with the classic bookinfo app.


### Testing with manual sidecar injection

Run the following commands:
```
oc new-project bookinfo
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f <(istioctl kube-inject --injectConfigMapName istio-sidecar-injector -f https://raw.githubusercontent.com/istio/istio/release-0.8/samples/bookinfo/kube/bookinfo.yaml) -n bookinfo
```

At this point the pods should deploy correclty and should participate to the istio mesh.

### Testing with automatic sidecar injection 

Run the following commands:
```
oc new-project bookinfo
oc adm policy add-scc-to-user anyuid -z default -n bookinfo
oc apply -f https://raw.githubusercontent.com/istio/istio/release-0.8/samples/bookinfo/kube/bookinfo.yaml -n bookinfo
```

At this point the pods should deploy correclty and should participate to the istio mesh.

## Installation with the applier

The applier is a infrastrcuture as code approach used in the OpenShift Community of practice.
Clone Repository

```
mkdir -p $GOPATH/src/github.com/sabre1041
cd $GOPATH/src/github.com/sabre1041
git clone https://github.com/sabre1041/istio-pod-network-controller.git
cd istio-pod-network-controller
```

Login to OpenShift Environment

```
oc login https://<MASTER_API>
```

Run Ansible Galaxy to retrieve dependencies

```
ansible-galaxy install -r requirements.yml --roles-path=galaxy
```

Deployment

```
ansible-playbook -i ./inventory galaxy/openshift-applier/playbooks/openshift-cluster-seed.yml
```

## Building

To build the code install `dep`:
```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

initialize the dependecies:
```
dep ensure -vendor-only
```
then build the code:
```
go build -v -o bin/istio-pod-network-controller main.go
```





