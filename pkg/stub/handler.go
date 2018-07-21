package stub

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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
	out, err := exec.Command("docker ps | grep " + pod.ObjectMeta.podID + " | grep k8s_POD | awk '{print $1}'").Output()
	if err != nil {
		logrus.Errorf("Failed to get containerID : %v", err)
		return
	}
	containerID := out.String()
	out, err := exec.Command("docker inspect --format {{.State.Pid}} " + containerID).Output()
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return
	}
	pidID := out.String()
	out, err := exec.Command("nsenter -t " + pidID + " -n istio-iptables.sh <params>").Output()
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return
	}
	return nil
}
