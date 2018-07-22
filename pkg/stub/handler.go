package stub

import (
	"context"

	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"os/exec"
	"sort"
)

func NewHandler(nodeName string) sdk.Handler {
	return &Handler{nodeName: nodeName}
}

type Handler struct {
	nodeName string
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *corev1.Pod:

		// Check to see if pod is running on current node
		if h.nodeName == o.Spec.NodeName {
			if filterPod(o) {
				err := managePod(o)
				if err != nil {
					logrus.Errorf("Failed to process pod : %v", err)
					return err
				}
			}
		}
	}
	return nil
}

func filterPod(pod *corev1.Pod) bool {
	// filter by state
	if pod.Status.Phase != "Running" && pod.Status.Phase != "Pending" {
		logrus.Infof("Pod %s terminated, ignoring", pod.ObjectMeta.Name)
		return false
	}
	// filter by whether the pod belongs to the mesh
	// not sure how to do it right now.

	//filter by opted in namespaces
	namespaces := []string{"tutorial"}
	sort.Strings(namespaces)
	i := sort.SearchStrings(namespaces, pod.ObjectMeta.Namespace)
	if !(i < len(namespaces) && namespaces[i] == pod.ObjectMeta.Namespace) {
		logrus.Infof("Pod %s not in considered namespaces, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// filter by being already initialized
	if "true" == pod.ObjectMeta.Annotations["initializer.istio.io/status"] {
		logrus.Infof("Pod %s previously initialized, ignoring", pod.ObjectMeta.Name)
		return false
	}
	return true
}

func managePod(pod *corev1.Pod) error {
	logrus.Infof("Processing Pod: %s", pod.ObjectMeta.Name)
	cmd := "-c docker ps | grep " + string(pod.ObjectMeta.UID) + " | grep k8s_POD | awk '{print $1}'"
	out, err := exec.Command("/bin/bash", cmd).Output()
	if err != nil {
		logrus.Errorf("Failed to get containerID : %v", err)
		return err
	}
	containerID := fmt.Sprintf("%s", out)
	logrus.Infof("ose_pod container id: %s", containerID)
	out, err = exec.Command("/bin/bash", "-c docker inspect --format {{.State.Pid}} "+containerID).Output()
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return err
	}
	pidID := fmt.Sprintf("%s", out)
	logrus.Infof("ose_pod container main process id: %s", pidID)
	out, err = exec.Command("/bin/bash", "-c nsenter -t "+pidID+" -n /usr/local/bin/istio-iptables.sh $ISTIO_PARAMS").Output()
	if err != nil {
		logrus.Errorf("Failed to setup ip tables : %v", err)
		return err
	}
	logrus.Infof("ip tables updated with no error")
	updatedPod := pod.DeepCopy()
	updatedPod.ObjectMeta.Annotations["initializer.istio.io/status"] = "true"
	err = sdk.Update(updatedPod)
	return err
}
