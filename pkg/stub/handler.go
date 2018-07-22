package stub

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

var defaultTimeout = 10 * time.Second

func NewHandler(nodeName string, dockerClient client.Client) sdk.Handler {
	return &Handler{nodeName: nodeName, dockerClient: dockerClient}
}

type Handler struct {
	nodeName     string
	dockerClient client.Client
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *corev1.Pod:

		// Check to see if pod is running on current node
		if h.nodeName == o.Spec.NodeName {

			err := managePod(h, ctx, o)
			if err != nil {
				logrus.Errorf("Failed to process pod : %v", err)
				return err
			}
		}

	}
	return nil
}

func managePod(h *Handler, ctx context.Context, pod *corev1.Pod) error {

	podName := pod.ObjectMeta.Name
	podNamespace := pod.ObjectMeta.Namespace

	dockerCtx, cancelFn := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFn()

	filter := filters.NewArgs()
	filter.Add("label", "io.kubernetes.container.name=POD")
	filter.Add("label", fmt.Sprintf("io.kubernetes.pod.name=%s", podName))
	filter.Add("label", fmt.Sprintf("io.kubernetes.pod.namespace=%s", podNamespace))

	containers, err := h.dockerClient.ContainerList(dockerCtx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(containers) == 1 {
		inspect, err := h.dockerClient.ContainerInspect(dockerCtx, containers[0].ID)

		if err != nil {
			logrus.Error(err)
			return err
		}

		logrus.Infof("Pod Namespace: %s - Pod Name: %s - Container PID: %d", podName, podNamespace, inspect.State.Pid)
	}

	return nil
}
