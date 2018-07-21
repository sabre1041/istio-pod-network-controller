package stub

import (
	"context"

	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"os/exec"
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

			err := managePod(o)
			if err != nil {
				logrus.Errorf("Failed to process pod : %v", err)
				return err
			}
		}

	}
	return nil
}

func managePod(pod *corev1.Pod) error {
	logrus.Infof("Processing Pod: %s", pod.ObjectMeta.Name)
	cmd := "docker ps | grep " + string(pod.ObjectMeta.UID) + " | grep k8s_POD | awk '{print $1}'"
	out, err := exec.Command(cmd).Output()
	if err != nil {
		logrus.Errorf("Failed to get containerID : %v", err)
		return err
	}
	containerID := fmt.Sprintf("%s", out)
	logrus.Infof("ose_pod container id: %s", containerID)
	out, err = exec.Command("docker inspect --format {{.State.Pid}} " + containerID).Output()
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return err
	}
	pidID := fmt.Sprintf("%s", out)
	logrus.Infof("ose_pod container main process id: %s", pidID)
	out, err = exec.Command("nsenter -t " + pidID + " -n /usr/local/bin/istio-iptables.sh $ISTIO_PARAMS").Output()
	if err != nil {
		logrus.Errorf("Failed to setup ip tables : %v", err)
		return err
	}
	logrus.Infof("ip tables update with no error")
	return nil
}
