package getpodmetric

import (
	"context"
	"fmt"
	"pro/common"
	"pro/listallnamespace"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetPodMetric() map[string]string {
	PodMetricMap := make(map[string]string)
	config, err := common.GetConfig()
	if err != nil {
		klog.Fatalf("Failed to create client config: %v", err)
	}
	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespaces := listallnamespace.ListAllNamespce()
	for _, ns := range namespaces {
		if ns == "kube-system" || ns == "kube-public" || ns == "kube-node-lease" || ns == "istio-system" || ns == "kube-flannel" {
			continue
		}
		podMetrics, err := mc.MetricsV1beta1().PodMetricses(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Println("Error:", err)
		}
		for _, podMetric := range podMetrics.Items {
			podName := podMetric.Name
			podName_ := podName
			// 使用 '-' 作为分隔符拆分字符串
			parts := strings.Split(podName, "-")
			// 获取前两个子字符串
			if len(parts) > 2 {
				podName_ = strings.Join(parts[:2], "-")
			} else {
				fmt.Println("字符串格式不正确")
			}
			podContainers := podMetric.Containers
			cpuQuantity := ""
			memQuantity := ""
			for _, container := range podContainers {
				if container.Name == "istio-proxy" {
					continue
				}
				cpuQuantity = container.Usage.Cpu().String()
				memQuantity = container.Usage.Memory().String()
				PodMetricMap["/pod"+"/"+ns+"/"+podName_+"/"+"cpuUsage"] = cpuQuantity
				PodMetricMap["/pod"+"/"+ns+"/"+podName_+"/"+"memUsage"] = memQuantity
			}

		}
	}
	return PodMetricMap
}
