package stub

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var sortedNamespaces []string

var defaultTimeout = 10 * time.Second

func NewHandler(nodeName string, dockerClient client.Client) sdk.Handler {
	sortedNamespaces = strings.Split(",", os.Getenv("NAMESPACES"))
	sort.Strings(sortedNamespaces)
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
			if filterPod(o) {
				err := managePod(h, ctx, o)
				if err != nil {
					logrus.Errorf("Failed to process pod : %v", err)
					return err
				}
				err = markPodAsInitialized(o)
				if err != nil {
					logrus.Errorf("Failed to process pod : %v", err)
					return err
				}
			}
		}
	}
	return nil
}

func markPodAsInitialized(pod *corev1.Pod) error {
	updatedPod := pod.DeepCopy()
	updatedPod.ObjectMeta.Annotations["pod-network-controller.istio.io/status"] = "initialized"
	err := sdk.Update(updatedPod)
	return err
}

func getPid(h *Handler, ctx context.Context, pod *corev1.Pod) (string, error) {

	dockerCtx, cancelFn := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFn()

	podName := pod.ObjectMeta.Name
	podNamespace := pod.ObjectMeta.Namespace
	filter := filters.NewArgs()
	filter.Add("label", "io.kubernetes.container.name=POD")
	filter.Add("label", fmt.Sprintf("io.kubernetes.pod.name=%s", podName))
	filter.Add("label", fmt.Sprintf("io.kubernetes.pod.namespace=%s", podNamespace))

	containers, err := h.dockerClient.ContainerList(dockerCtx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if len(containers) == 1 {
		inspect, err := h.dockerClient.ContainerInspect(dockerCtx, containers[0].ID)

		if err != nil {
			logrus.Error(err)
			return "", err
		}

		logrus.Infof("Pod Namespace: %s - Pod Name: %s - Container PID: %d", podName, podNamespace, inspect.State.Pid)
		return fmt.Sprintf("%d", inspect.State.Pid), nil
	}
	return "", errors.New("unable to find pod main pid")
}

func filterPod(pod *corev1.Pod) bool {
	// filter by state
	if pod.Status.Phase != "Running" && pod.Status.Phase != "Pending" {
		logrus.Infof("Pod %s terminated, ignoring", pod.ObjectMeta.Name)
		return false
	}
	// filter by whether the pod belongs to the mesh
	// not sure how to do it right now.

	// make sure the pod if not a deployer pod
	if _, ok := pod.ObjectMeta.Labels["openshift.io/deployer-pod-for.name"]; ok {
		return false
	}

	// make sure the pod if not a build pod
	if _, ok := pod.ObjectMeta.Labels["openshift.io/build.name"]; ok {
		return false
	}

	//filter by opted in namespaces
	//	namespaces := []string{"tutorial"}
	//	sort.Strings(namespaces)
	i := sort.SearchStrings(sortedNamespaces, pod.ObjectMeta.Namespace)
	if !(i < len(sortedNamespaces) && sortedNamespaces[i] == pod.ObjectMeta.Namespace) {
		logrus.Infof("Pod %s not in considered namespaces, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// filter by being already initialized
	if "initialized" == pod.ObjectMeta.Annotations["pod-network-controller.istio.io/status"] {
		logrus.Infof("Pod %s previously initialized, ignoring", pod.ObjectMeta.Name)
		return false
	}
	return true
}

func managePod(h *Handler, ctx context.Context, pod *corev1.Pod) error {
	logrus.Infof("Processing Pod: %s , id %s", pod.ObjectMeta.Name, pod.ObjectMeta.UID)
	pidID, err := getPid(h, ctx, pod)
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return err
	}
	logrus.Infof("ose_pod container main process id: %s", pidID)
	args1 := []string{"-t", pidID, "-n", "/usr/local/bin/istio-iptables.sh"}
	out, err := exec.Command("nsenter", append(args1, os.Args[1:]...)...).CombinedOutput()
	logrus.Infof("nsenter output: %s", out)
	if err != nil {
		logrus.Errorf("Failed to setup ip tables : %v", err)
		return err
	}
	logrus.Infof("ip tables updated with no error")
	return err
}
