package run

import (
	"context"
	"github.com/docker/docker/client"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

var log = logrus.New()

func initLog() {
	var err error
	log.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalln(err)
	}
}

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func Run() {
	initLog()
	//out, _ := exec.Command("docker", "ps").CombinedOutput()
	//logrus.Infof("%s", out)
	printVersion()

	nodeName := os.Getenv("NODE_NAME")

	if nodeName == "" {
		logrus.Error("NODE_NAME not defined")
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	logrus.Infof("Managing Pods Running on Node: %s", nodeName)
	sdk.Watch("v1", "Pod", "", 0)
	sdk.Handle(NewHandler(nodeName, *cli))
	sdk.Run(context.TODO())
}
