package main

import (
	"context"
	"os"
	"os/exec"
	"runtime"

	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	stub "github.com/sabre1041/istio-pod-network-controller/pkg/stub"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	//out, _ := exec.Command("docker", "ps").CombinedOutput()
	//logrus.Infof("%s", out)
	printVersion()

	nodeName := os.Getenv("NODE_NAME")

	if nodeName == "" {
		logrus.Error("NODE_NAME not defined")
		return
	}

	logrus.Infof("Managing Pods Running on Node: %s", nodeName)
	logrus.Infof("Raffa - 1")
	sdk.Watch("v1", "Pod", "", 0)
	sdk.Handle(stub.NewHandler(nodeName))
	sdk.Run(context.TODO())
}
