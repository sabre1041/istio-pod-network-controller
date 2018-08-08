package init

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const PodAnnotationsFileName = "/etc/podinfo/pod_annotations"
const PodAnnotationsKeyName = "pod-network-controller.istio.io/status"
const PodAnnotationsValueName = "initialized"

func WaitForAnnotationInFile(filePath string, annotationKey string, annotationValue string) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("File '%s' Does Not Exist", filePath)
	}

	for {
		result, err := checkForAnnotation(filePath, annotationKey, annotationValue)

		if result {
			break
		} else if err != nil {
			return err
		}

		// Delay for 5 seconds
		time.Sleep(time.Second * 5)

	}

	return nil
}

func checkForAnnotation(filePath string, annotationKey string, annotationValue string) (bool, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return false, fmt.Errorf("Error accessing file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

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
