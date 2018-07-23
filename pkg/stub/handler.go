package stub

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	//	"os/exec"
	"errors"
	"sort"
	"time"
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
			if filterPod(o) {
				err := managePod(h, ctx, o)
				if err != nil {
					logrus.Errorf("Failed to process pod : %v", err)
					return err
				}
			}
		}
	}
	return nil
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
		return fmt.Sprintf("%s", inspect.State.Pid), nil
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

func managePod(h *Handler, ctx context.Context, pod *corev1.Pod) error {
	logrus.Infof("Processing Pod: %s , id %s", pod.ObjectMeta.Name, pod.ObjectMeta.UID)
	//	cmd := "-c docker ps | grep " + string(pod.ObjectMeta.UID)
	//	out, err := exec.Command("/bin/bash", cmd).Output()
	//	if err != nil {
	//		logrus.Errorf("Failed to get containerID : %v", err)
	//		return err
	//	}
	//	logrus.Infof("output command 1: %s", out)
	//	cmd = "-c docker ps | grep " + string(pod.ObjectMeta.UID) + " | grep k8s_POD "
	//	out, err = exec.Command("/bin/bash", cmd).Output()
	//	if err != nil {
	//		logrus.Errorf("Failed to get containerID : %v", err)
	//		return err
	//	}
	//	logrus.Infof("output command 2: %s", out)
	//cmd := "-c docker ps | grep " + fmt.Sprintf("%s", pod.ObjectMeta.UID) + " | grep k8s_POD | awk '{print $1}'"
	//out, err := exec.Command("/bin/bash", cmd).Output()

	//	out, err := exec.Command("docker", "ps", "--filter", "label=io.kubernetes.container.name=POD", "--filter", "label=io.kubernetes.pod.name="+pod.ObjectMeta.Name, "-q").CombinedOutput()
	//	//out, err := exec.Command("/bin/bash", "-c", "docker", "ps", "|", "grep", pod.ObjectMeta.Name, "|", "grep", "k8s_POD", "|", "awk", "'{print $1}'").CombinedOutput()
	//	if err != nil {
	//		logrus.Errorf("Failed to get containerID : %v", err)
	//		return err
	//	}
	//	containerID := fmt.Sprintf("%s", out)
	//	logrus.Infof("ose_pod container id: %s", containerID)
	//	out, err = exec.Command("docker", "inspect", "--format", "{{.State.Pid}}", containerID).CombinedOutput()
	//	if err != nil {
	//		logrus.Errorf("Failed to get pidID : %v", err)
	//		return err
	//	}
	pidID, err := getPid(h, ctx, pod)
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return err
	}
	logrus.Infof("ose_pod container main process id: %s", pidID)
	//	out, err = exec.Command("nsenter", "-t", pidID, "-n", "/usr/local/bin/istio-iptables.sh", "$ISTIO_PARAMS").CombinedOutput()
	//	if err != nil {
	//		logrus.Errorf("Failed to setup ip tables : %v", err)
	//		return err
	//	}
	//	logrus.Infof("ip tables updated with no error")
	//	updatedPod := pod.DeepCopy()
	//	updatedPod.ObjectMeta.Annotations["initializer.istio.io/status"] = "true"
	//	err = sdk.Update(updatedPod)
	return err
}
