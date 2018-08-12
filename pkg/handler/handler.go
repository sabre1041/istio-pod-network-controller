package run

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const EnvoyProxyUserID = "1337"
const EnvoyPortAnnotation = "pod-network-controller.istio.io/envoy-port"
const InterceptModeAnnotation = "sidecar.istio.io/interceptionMode"
const IncludePortsAnnotation = "traffic.sidecar.istio.io/includeInboundPorts"
const ExcludePortsAnnotation = "traffic.sidecar.istio.io/excludeInboundPorts"
const IncludeCidrsAnnotation = "traffic.sidecar.istio.io/includeOutboundIPRanges"
const ExcludeCidrsAnnotation = "traffic.sidecar.istio.io/excludeOutboundIPRanges"
const EnvoyUseridAnnotation = "pod-network-controller.istio.io/envoy-userid"
const EnvoyGroupidAnnotation = "pod-network-controller.istio.io/envoy-groupid"
const TargetedPodAnnotation = "istio-pod-network-controller/initialize"
const TargetedPodAnnotationValue = "true"
const PodNetworkControllerAnnotation = "pod-network-controller.istio.io/status"
const PodNetworkControllerAnnotationInitialized = "initialized"
const DeployerPodAnnotation = "openshift.io/deployer-pod-for.name"
const BuildPodAnnotation = "openshift.io/build.name"

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
	updatedPod.ObjectMeta.Annotations[PodNetworkControllerAnnotation] = PodNetworkControllerAnnotationInitialized
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
		logrus.Debugf("Pod %s terminated, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// make sure the pod if not a deployer pod
	if _, ok := pod.ObjectMeta.Labels[DeployerPodAnnotation]; ok {
		logrus.Debugf("Pod %s is a deployer, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// make sure the pod if not a build pod
	if _, ok := pod.ObjectMeta.Labels[BuildPodAnnotation]; ok {
		logrus.Debugf("Pod %s is a builder, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// filter by being already initialized
	if PodNetworkControllerAnnotationInitialized == pod.ObjectMeta.Annotations[PodNetworkControllerAnnotation] {
		logrus.Infof("Pod %s previously initialized, ignoring", pod.ObjectMeta.Name)
		return false
	}

	// Check if Pod Annotated for Injection
	if TargetedPodAnnotationValue == pod.ObjectMeta.Annotations[TargetedPodAnnotation] {
		return true
	}

	// Check if Namespace Annotated for Injection
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.ObjectMeta.Namespace,
		},
	}

	err := sdk.Get(namespace)

	if err != nil {
		logrus.Errorf("Failed to namespace for pod : %v", err)
	}

	if TargetedPodAnnotationValue == namespace.ObjectMeta.Annotations[TargetedPodAnnotation] {
		return true
	}

	return false
}

func managePod(h *Handler, ctx context.Context, pod *corev1.Pod) error {
	logrus.Infof("Processing Pod: %s , id %s", pod.ObjectMeta.Name, pod.ObjectMeta.UID)
	pidID, err := getPid(h, ctx, pod)
	if err != nil {
		logrus.Errorf("Failed to get pidID : %v", err)
		return err
	}
	logrus.Infof("ose_pod container main process id: %s", pidID)
	args := []string{"-t", pidID, "-n", "/usr/local/bin/istio-iptables.sh", "-p", getEnvoyPort(pod),
		"-u", getUserID(pod), "-g", getGroupID(pod), "-m", getInterceptMode(pod), "-b", getIncludedInboundPorts(pod), "-d", getExcludedInboundPorts(pod),
		"-i", getIncludedOutboundCidrs(pod), "-x", getExcludedOutboundCidrs(pod)}
	logrus.Infof("excuting ip tables rules with the following arguments: %s", args)
	out, err := exec.Command("nsenter", args...).CombinedOutput()
	logrus.Infof("nsenter output: %s", out)
	if err != nil {
		logrus.Errorf("Failed to setup ip tables : %v", err)
		return err
	}
	logrus.Infof("ip tables updated with no error")
	return err
}

func getEnvoyPort(pod *corev1.Pod) string {
	if port, ok := pod.ObjectMeta.Labels[EnvoyPortAnnotation]; ok {
		return port
	} else {
		return viper.GetString("envoy-port")
	}
}

func getUserID(pod *corev1.Pod) string {
	return EnvoyProxyUserID
}

func getGroupID(pod *corev1.Pod) string {
	return EnvoyProxyUserID
}

func getInterceptMode(pod *corev1.Pod) string {
	if interceptMode, ok := pod.ObjectMeta.Labels[InterceptModeAnnotation]; ok {
		return interceptMode
	} else {
		return viper.GetString("istio-inbound-interception-mode")
	}
}

func getIncludedInboundPorts(pod *corev1.Pod) string {
	if includePorts, ok := pod.ObjectMeta.Labels[IncludePortsAnnotation]; ok {
		return includePorts
	} else {
		ports := ""
		for _, k := range pod.Spec.Containers {
			for _, p := range k.Ports {
				ports += fmt.Sprintf("%d", p.ContainerPort) + ","
			}
		}
		return ports
	}
}

func getPodServicePorts(pod *corev1.Pod) string {
	return ""
}

func getExcludedInboundPorts(pod *corev1.Pod) string {
	if port, ok := pod.ObjectMeta.Labels[ExcludePortsAnnotation]; ok {
		return port
	} else {
		return ""
	}
}

func getIncludedOutboundCidrs(pod *corev1.Pod) string {
	if includeCidrs, ok := pod.ObjectMeta.Labels[IncludeCidrsAnnotation]; ok {
		return includeCidrs
	} else {
		return "*"
	}
}

func getExcludedOutboundCidrs(pod *corev1.Pod) string {
	if excludeCidrs, ok := pod.ObjectMeta.Labels[ExcludeCidrsAnnotation]; ok {
		return excludeCidrs
	} else {
		return ""
	}
}
