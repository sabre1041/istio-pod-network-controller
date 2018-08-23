# Build instructions


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