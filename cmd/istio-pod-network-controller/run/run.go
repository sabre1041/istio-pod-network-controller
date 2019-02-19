// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package run

import (
	"context"

	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/docker/docker/client"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	handler "github.com/sabre1041/istio-pod-network-controller/pkg/handler"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func initLog() {
	var err error
	log.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalln(err)
	}
}

func NewRunCmd() *cobra.Command {

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "starts the istio pod network controller",
		Long:  "starts the istio pod network controller",
		Run:   runFunc,
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().String("enable-inbound-ipv6", "false", "whether inbound ipv6 connection should be managed by the mesh (currently is must be set to false)")
	runCmd.Flags().String("envoy-port", "15001", "Specify the envoy port to which redirect all TCP traffic. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+handler.EnvoyPortAnnotation)
	runCmd.Flags().String("istio-inbound-interception-mode", "REDIRECT", "The mode used to redirect inbound connections to Envoy, either REDIRECT or TPROXY. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+handler.InterceptModeAnnotation)
	runCmd.Flags().String("istio-include-outbound-ip-ranges", "*", "The istio egress capture IP whitelist. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+handler.IncludeCidrsAnnotation)
	runCmd.Flags().String("istio-exclude-outbound-ip-ranges", "", "The istio egress capture IP exception list. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+handler.ExcludeCidrsAnnotation)
	runCmd.Flags().String("istio-exclude-inbound-ports", "", "The istio ingress capture port exception list. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+handler.ExcludePortsAnnotation)
	runCmd.Flags().String("container-runtime", "docker", "container runtime, suppported values are: 'docker' and 'crio'")
	runCmd.Flags().String("crio-socket", "unix:///var/run/crio/crio.sock", "the socker where the cri server is at")
	runCmd.Flags().String("runc-root", "/run/runc", "root directory for runc's storage of container state")
	runCmd.Flags().String("node-name", "", "the node that should be monitored, pass this with the downward API")
	viper.BindPFlag("enable-inbound-ipv6", runCmd.Flags().Lookup("enable-inbound-ipv6"))
	viper.BindPFlag("envoy-port", runCmd.Flags().Lookup("envoy-port"))
	viper.BindPFlag("istio-inbound-interception-mode", runCmd.Flags().Lookup("istio-inbound-interception-mode"))
	viper.BindPFlag("istio-include-outbound-ip-ranges", runCmd.Flags().Lookup("istio-include-outbound-ip-ranges"))
	viper.BindPFlag("istio-exclude-outbound-ip-ranges", runCmd.Flags().Lookup("istio-exclude-outbound-ip-ranges"))
	viper.BindPFlag("istio-exclude-inbound-ports", runCmd.Flags().Lookup("istio-exclude-inbound-ports"))
	viper.BindPFlag("container-runtime", runCmd.Flags().Lookup("container-runtime"))
	viper.BindPFlag("node-name", runCmd.Flags().Lookup("node-name"))
	viper.BindPFlag("crio-socket", runCmd.Flags().Lookup("crio-socket"))
	viper.BindPFlag("runc-root", runCmd.Flags().Lookup("runc-root"))

	return runCmd

}

func printVersion() {
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func runFunc(cmd *cobra.Command, args []string) {
	initLog()

	printVersion()

	/*	var containerRuntime string

		if "crio" == viper.GetString("container-runtime") {
			out, err := exec.Command("/bin/bash", "-c", "crio config | grep \"^runtime .*\" | awk '{print $3}' | tr -d '\"' ").CombinedOutput()
			log.Infof("container runtime output: %s", out)
			if err != nil {
				log.Error("couldn't retrieve container runtime executable: %s", err)
				return
			}
			containerRuntime = fmt.Sprintf("%s", out)
		}*/

	if "" == viper.GetString("node-name") {
		log.Error("NODE_NAME not defined")
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	log.Infof("Managing Pods Running on Node: %s", viper.GetString("node-name"))
	sdk.Watch("v1", "Pod", "", 0)
	sdk.Handle(handler.NewHandler(viper.GetString("node-name"), *cli, "runc"))
	sdk.Run(context.TODO())
}
