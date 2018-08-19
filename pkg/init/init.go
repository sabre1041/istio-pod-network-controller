package init

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const PodAnnotationsFileName = "/etc/podinfo/pod_annotations"
const PodAnnotationsKeyName = "pod-network-controller.istio.io/status"
const PodAnnotationsValueName = "initialized"
const InitTimeout = 300
const InitDelay = 10

func WaitForAnnotationInFile(filePath string, annotationKey string, annotationValue string, timeout time.Duration, delay int) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("File '%s' Does Not Exist", filePath)
	}

	done := make(chan interface{}, 1)
	var resultError error

	deadline := time.Now().Add(timeout)

	go func() {
		defer close(done)

		for time.Now().Before(deadline) {
			result, err := checkForAnnotation(filePath, annotationKey, annotationValue)

			if result {
				return
			} else if err != nil {
				resultError = err
				return
			}

			time.Sleep(time.Duration(delay) * time.Second)

		}
	}()

	select {
	case <-time.After(timeout):
		return errors.New("Timed out waiting for pod annotation")
	case <-done:
		return resultError
	}
}

func checkForAnnotation(filePath string, annotationKey string, annotationValue string) (bool, error) {

	fileContents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return false, fmt.Errorf("Error accessing file: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(fileContents)))

	for scanner.Scan() {
		line := scanner.Text()

		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}

				if annotationKey == trimQuotes(key) && annotationValue == trimQuotes(value) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
