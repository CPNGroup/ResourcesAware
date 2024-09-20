package getpodmetric

import (
	"context"
	"fmt"
	"pro/common"
	"pro/listallnamespace"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

func GetPodResources() map[string]string {
	PodMetricMap := make(map[string]string)

	// 获取 Kubernetes 配置
	config, err := common.GetConfig()
	if err != nil {
		klog.Fatalf("Failed to create client config: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 获取所有命名空间
	namespaces := listallnamespace.ListAllNamespce()

	for _, ns := range namespaces {
		// if ns == "kube-system" || ns == "kube-public" || ns == "kube-node-lease" || ns == "istio-system" || ns == "kube-flannel" {
		// 	continue
		// }
		// 这里只列出我自己的命名空间
		if ns == "li" {
			// 列出命名空间下的所有 Pod
			podList, err := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				fmt.Println("Error listing pods:", err)
				continue
			}

			for _, pod := range podList.Items {
				podName := pod.Name
				podName_ := podName
				// 使用 '-' 作为分隔符拆分字符串
				parts := strings.Split(podName, "-")
				// 获取前两个子字符串
				if len(parts) > 2 {
					podName_ = strings.Join(parts[:2], "-")
				} else {
					fmt.Println("字符串格式不正确")
				}
				for _, container := range pod.Spec.Containers {
					if container.Name == "istio-proxy" {
						continue
					}

					// 获取容器的资源请求和限制
					cpuRequest := container.Resources.Requests[corev1.ResourceCPU]
					memRequest := container.Resources.Requests[corev1.ResourceMemory]
					cpuLimit := container.Resources.Limits[corev1.ResourceCPU]
					memLimit := container.Resources.Limits[corev1.ResourceMemory]

					// // 计算容器的资源总量
					// cpuQuantity := cpuRequest.String() + "/" + cpuLimit.String()
					// memQuantity := memRequest.String() + "/" + memLimit.String()

					// 将容器资源总量存入 map 中
					PodMetricMap["/pod/"+ns+"/"+podName_+"/cpuRequest"] = cpuRequest.String()
					PodMetricMap["/pod/"+ns+"/"+podName_+"/memRequest"] = memRequest.String()
					PodMetricMap["/pod/"+ns+"/"+podName_+"/cpuLimit"] = cpuLimit.String()
					PodMetricMap["/pod/"+ns+"/"+podName_+"/memLimit"] = memLimit.String()
				}
			}
		}
	}
	return PodMetricMap
}
