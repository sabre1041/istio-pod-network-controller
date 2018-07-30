#!/usr/bin/env bash
docker build -t quay.io/raffaelespazzoli/istio-pod-network-controller:latest .
docker push quay.io/raffaelespazzoli/istio-pod-network-controller:latest