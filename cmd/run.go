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

package cmd

import (
	run "github.com/sabre1041/istio-pod-network-controller/cmd/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the istio pod network controller",
	Long:  "starts the istio pod network controller",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		run.Run()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().String("include-namespaces", "", "comma-separated list of namespace that should be considered (ex: namespace1,namespace2,namespace3).")
	runCmd.Flags().String("include-namespaces-regexp", "", "regular expression that identifies manespaces that shoul be considered ((not implemented yet).")
	runCmd.Flags().String("exclude-namespaces", "", "comma-separated list of namespaces that should not be considered (ex: namespace1,namespace2,namespace3).")
	runCmd.Flags().String("exclude-namespaces-regexp", "", "regular expression that identifies namespaces that should not be considered (not implemented yet).")
	runCmd.Flags().String("enable-inbound-ipv6", "false", "whether inbound ipv6 connection should be managed by the mesh (currently is must be set to false)")
	runCmd.Flags().String("envoy-port", "15001", "Specify the envoy port to which redirect all TCP traffic. This is a cluster-wide setting, you can override it by adding this annotation to your pods "+run.EnvoyPortAnnotation)
	runCmd.Flags().String("istio-inbound-interception-mode", "REDIRECT", "The mode used to redirect inbound connections to Envoy, either REDIRECT or TPROXY. his is a cluster-wide setting, you can override it by adding this annotation to your pods "+run.InterceptModeAnnotation)
	//runCmd.Flags().String("istio-include-inbound-ports", "", "comma-separated list of ports that will be redirected to Envoy. This list will be added to the ports desumed from the service controlling the pod. The wildcard character \"*\" can be used to configure redirection for all ports. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.IncludePortsAnnotation)
	//runCmd.Flags().String("istio-exclude-inbound-ports", "", "Comma separated list of inbound ports to be excluded from redirection to Envoy (optional). Only applies when all inbound traffic (i.e. \"*\") is being redirected. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.ExcludePortsAnnotation)
	//runCmd.Flags().String("istio-include-outbound-cidrs", "*", "Comma separated list of IP ranges in CIDR form to redirect to envoy (optional). The wildcardcharacter \"*\" can be used to redirect all outbound traffic. An empty list will disable all outbound redirection. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.IncludeCidrsAnnotation)
	//runCmd.Flags().String("istio-exclude-outbound-cidrs", "", "Comma separated list of IP ranges in CIDR form to be excluded from redirection. Only applies when all outbound traffic (i.e. \"*\") is being redirected. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.ExcludeCidrsAnnotation)
	//runCmd.Flags().String("envoy-userid", "1337", "UID used by the envoy-proxy container. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.EnvoyUseridAnnotation)
	//runCmd.Flags().String("envoy-groupid", "1337", "GID used by the envoy-proxy container. This is a cluster-wide settings, you can override it by adding this annotation to your pods "+run.EnvoyGroupidAnnotation)
	viper.BindPFlag("include-namespaces", runCmd.Flags().Lookup("include-namespaces"))
	viper.BindPFlag("include-namespaces-regexp", runCmd.Flags().Lookup("include-namespaces-regexp"))
	viper.BindPFlag("exclude-namespaces", runCmd.Flags().Lookup("exclude-namespaces"))
	viper.BindPFlag("exclude-namespaces-regexp", runCmd.Flags().Lookup("exclude-namespaces-regexp"))
	viper.BindPFlag("enable-inbound-ipv6", runCmd.Flags().Lookup("enable-inbound-ipv6"))
	viper.BindPFlag("envoy-port", runCmd.Flags().Lookup("envoy-port"))
	viper.BindPFlag("istio-inbound-interception-mode", runCmd.Flags().Lookup("istio-inbound-interception-mode"))
	//viper.BindPFlag("istio-include-inbound-ports", runCmd.Flags().Lookup("istio-include-inbound-ports"))
	//viper.BindPFlag("istio-exclude-inbound-ports", runCmd.Flags().Lookup("istio-exclude-inbound-ports"))
	//viper.BindPFlag("istio-include-outbound-cidrs", runCmd.Flags().Lookup("istio-include-outbound-cidrs"))
	//viper.BindPFlag("istio-exclude-outbound-cidrs", runCmd.Flags().Lookup("istio-exclude-outbound-cidrs"))
	//viper.BindPFlag("envoy-userid", runCmd.Flags().Lookup("envoy-userid"))
	//viper.BindPFlag("envoy-groupid", runCmd.Flags().Lookup("envoy-groupid"))

}
