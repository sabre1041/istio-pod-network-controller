.DEFAULT_GOAL	:= build

#------------------------------------------------------------------------------
# Variables
#------------------------------------------------------------------------------

.PHONY: push
push: docker
	@docker push quay.io/raffaelespazzoli/istio-pod-network-controller:latest

.PHONY: docker
docker: build
	@echo "--> building docker image"
	@docker build -f Dockerfile -t quay.io/raffaelespazzoli/istio-pod-network-controller:latest .

.PHONY: build
build: vendor
	@echo "---> building go binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/istio-pod-network-controller -v cmd/istio-pod-network-controller/main.go

.PHONY: clean
clean:
	@echo "--> cleaning compiled objects and binaries"
	@go clean

.PHONY: check
check: vet

.PHONY: vet
vet: tools.govet
	@echo "--> checking code correctness with 'go vet' tool"
	@go vet ./cmd/...

#------------------
#-- dependencies
#------------------
.PHONY: depend.update depend.install

depend.update: tools.dep
	@echo "--> updating dependencies from Gopkg.yaml"
	@dep ensure update

depend.install: tools.dep
	@echo "--> installing dependencies from Gopkg.lock "
	@dep ensure

vendor: tools.dep
	@echo "--> installing dependencies from Gopkg.lock "
	@dep ensure -vendor-only

$(BINDIR):
	@mkdir -p $(BINDIR)

#---------------
#-- tools
#---------------
.PHONY: tools tools.dep tools.goimports tools.golint tools.govet

tools: tools.dep tools.goimports tools.golint tools.govet

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing goimports"; \
		go get golang.org/x/tools/cmd/goimports; \
	fi

tools.govet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		echo "--> installing govet"; \
		go get golang.org/x/tools/cmd/vet; \
	fi

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get -u golang.org/x/lint/golint; \
	fi

tools.dep:
	@command -v dep >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing dep"; \
		curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh; \
	fi
