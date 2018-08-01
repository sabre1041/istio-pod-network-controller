# Build instructions

This project can be built using one of two approaches:

1. Using the [openshift-applier](https://github.com/redhat-cop/openshift-applier)
2. Manually building the project

## Installation with the applier

The applier is a Infrastructure as Code (IaC) approach from the Red Hat Communities of Practice (CoP). The automation performed by the tool will result in build and deployment of the `istio-pod-network-controller`.

Refer to the steps illustrated on the project [README](./README.md) on deploying Istio.


1. Clone the Repository

    ```
    mkdir -p $GOPATH/src/github.com/sabre1041
    cd $GOPATH/src/github.com/sabre1041
    git clone https://github.com/sabre1041/istio-pod-network-controller.git
    cd istio-pod-network-controller
    ```

2. Login to OpenShift Environment

    ```
    oc login https://<MASTER_API>
    ```

3. Run Ansible Galaxy to retrieve the required dependencies

    ```
    ansible-galaxy install -r requirements.yml --roles-path=galaxy
    ```

4. Execute the `openshift-applier` to apply the configurations to the cluster

    ```
    ansible-playbook -i ./applier/inventory galaxy/openshift-applier/playbooks/openshift-cluster-seed.yml
    ```

## Manually Building the Project 

1. To build the project, first install the `dep` tool:

    ```
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
    ```

2. Retrieve the dependencies:

    ```
    dep ensure -vendor-only
    ```

3. Build the project:

    ```
    go build -v -o istio-pod-network-controller cmd/istio-pod-network-controller/main.go
    ```