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

	return nil
}
