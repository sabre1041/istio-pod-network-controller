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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	initial "github.com/sabre1041/istio-pod-network-controller/cmd/istio-pod-network-controller/init"
	"github.com/sabre1041/istio-pod-network-controller/cmd/istio-pod-network-controller/run"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// read in environment variables that match
	viper.AutomaticEnv()
}

func newRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "istio-pod-network-controller",
		Short: "A controller that adds pods to the istio service mesh",
		Long:  "A controller that adds pods to the istio service mesh",
	}

	rootCmd.PersistentFlags().String("log-level", "info", "log level")
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.AddCommand(run.NewRunCmd())
	rootCmd.AddCommand(initial.NewInitCmd())

	return rootCmd

}
