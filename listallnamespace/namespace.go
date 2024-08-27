package listallnamespace

import (
	"context"
	"time"

	"pro/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

func ListAllNamespce() []string {

	clientset, err := common.GetClientSet()
	if err != nil {
		klog.Fatalf("Failed to create client config: %v", err)
	}

	namespaces := listNamspace(clientset)
	return namespaces
}

func listNamspace(clientset *kubernetes.Clientset) []string {
	retries := 5
	namespaces := []string{}
	for retries > 0 {
		namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			klog.Errorf("Failed to list namespaces: %v", err)
			retries--
			time.Sleep(5 * time.Second) // 等待一段时间后重试
			continue
		}

		for _, namespace := range namespaceList.Items {
			namespaces = append(namespaces, namespace.Name)
		}

		break
	}
	return namespaces
}
