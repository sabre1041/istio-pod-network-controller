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

package init

import (
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdkVersion "github.com/operator-framework/operator-sdk/version"
	initial "github.com/sabre1041/istio-pod-network-controller/pkg/init"
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

func NewInitCmd() *cobra.Command {

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Performs initialization for the Istio pod Network Controller",
		Long:  "Performs initialization for the Istio pod Network Controller",
		Run:   initFunc,
	}

	initCmd.Flags().String("file", initial.PodAnnotationsFileName, "Location of the file containing Pod Annotations")
	initCmd.Flags().String("annotation-key", initial.PodAnnotationsKeyName, "Name of the Annotation Key to Wait For")
	initCmd.Flags().String("annotation-value", initial.PodAnnotationsValueName, "Name of the Annotation Value to Wait For")
	initCmd.Flags().Int("timeout", initial.InitTimeout, "Timeout value waiting for pod annotation")
	initCmd.Flags().Int("delay", initial.InitDelay, "Amount of time between checking whether pod has been annotated")

	viper.BindPFlag("file", initCmd.Flags().Lookup("file"))
	viper.BindPFlag("annotation-key", initCmd.Flags().Lookup("annotation-key"))
	viper.BindPFlag("annotation-value", initCmd.Flags().Lookup("annotation-value"))
	viper.BindPFlag("timeout", initCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("delay", initCmd.Flags().Lookup("delay"))

	return initCmd

}

func printVersion() {
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func initFunc(cmd *cobra.Command, args []string) {
	initLog()

	file := viper.GetString("file")
	annotationKey := viper.GetString("annotation-key")
	annotationValue := viper.GetString("annotation-value")
	delay := viper.GetInt("delay")
	timeout := viper.GetInt("timeout")

	log.Printf("Waiting for Initialized Pod Annotation (%s=%s)", annotationKey, annotationValue)

	err := initial.WaitForAnnotationInFile(file, annotationKey, annotationValue, time.Duration(timeout)*time.Second, delay)

	if err != nil {
		log.Errorf("Error occurred waiting for pod annotation in file: %v", err)
		os.Exit(1)
	}

	log.Printf("Completed Initialization Successfully")

}
