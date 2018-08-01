# Build instructions

You can build this project in two ways: with the applier of manually.

## Installation with the applier

The applier is a infrastructure as code approach used in the OpenShift Community of practice.
This automation will build the project and deploy the istio-pod-network-controller.
Refere to the manual deployment steps to deploy istio and the examples.


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
ansible-playbook -i ./applier/inventory galaxy/openshift-applier/playbooks/openshift-cluster-seed.yml
```

## Building the code 

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
go build -v -o istio-pod-network-controller cmd/istio-pod-network-controller/main.go
```