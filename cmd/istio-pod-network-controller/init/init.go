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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdkVersion "github.com/operator-framework/operator-sdk/version"
	initial "github.com/sabre1041/istio-pod-network-controller/pkg/init"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func NewInitCmd() *cobra.Command {

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Performs initialization for the Istio pod Network Controller",
		Long:  "Performs initialization for the Istio pod Network Controller",
		Run:   initFunc,
	}

	initCmd.Flags().String("file", getDefaultFileNameProperty(), "Location of the file containing Pod Annotations")
	initCmd.Flags().String("annotation-key", getDefaultPodAnnotationKeyProperty(), "Name of the Annotation Key to Wait For")
	initCmd.Flags().String("annotation-value", getDefaultPodAnnotationValueProperty(), "Name of the Annotation Value to Wait For")

	viper.BindPFlag("file", initCmd.Flags().Lookup("file"))
	viper.BindPFlag("annotation-key", initCmd.Flags().Lookup("annotation-key"))
	viper.BindPFlag("annotation-value", initCmd.Flags().Lookup("annotation-value"))

	return initCmd

}

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

func initFunc(cmd *cobra.Command, args []string) {
	initLog()

	file := viper.GetString("file")
	annotationKey := viper.GetString("annotation-key")
	annotationValue := viper.GetString("annotation-value")

	logrus.Printf("Waiting for Initialized Pod Annotation (%s=%s)", annotationKey, annotationValue)

	err := initial.WaitForAnnotationInFile(file, annotationKey, annotationValue)

	if err != nil {
		logrus.Errorf("Error occurred waiting for pod annotation in file: %v", err)
		os.Exit(1)
	}

	logrus.Printf("Completed Initialization Successfully")

}

func getDefaultFileNameProperty() string {
	fileProperty := os.Getenv("FILE")

	if fileProperty == "" {
		fileProperty = initial.PodAnnotationsFileName
	}

	return fileProperty
}

func getDefaultPodAnnotationKeyProperty() string {
	annotationKeyProperty := os.Getenv("ANNOTATION_KEY")

	if annotationKeyProperty == "" {
		annotationKeyProperty = initial.PodAnnotationsKeyName
	}

	return annotationKeyProperty
}

func getDefaultPodAnnotationValueProperty() string {
	annotationValueProperty := os.Getenv("ANNOTATION_VALUE")

	if annotationValueProperty == "" {
		annotationValueProperty = initial.PodAnnotationsValueName
	}

	return annotationValueProperty
}
