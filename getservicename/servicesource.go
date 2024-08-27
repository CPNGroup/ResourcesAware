package getservicename

import (
	"context"
	"encoding/json"
	"log"
	"pro/common"
	"pro/listallnamespace"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func ListServices() map[string]string {
	clientset, err := common.GetClientSet()
	if err != nil {
		klog.Fatalf("Failed to create client config: %v", err)
	}
	ServiceMetricMap := make(map[string]string)
	namespaces := listallnamespace.ListAllNamespce()
	for _, ns := range namespaces {
		if ns == "kube-system" || ns == "kube-public" || ns == "kube-node-lease" || ns == "istio-system" || ns == "kube-flannel" {
			continue
		}
		var services []string
		retries := 5
		for retries > 0 {
			svcList, err := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				klog.Errorf("Failed to list services: %v", err)
				retries--
				time.Sleep(5 * time.Second) // 等待一段时间后重试
				continue
			}

			for _, svc := range svcList.Items {
				services = append(services, svc.Name)
			}
			break
		}
		jsonServices, err := json.Marshal(services)
		if err != nil {
			log.Fatal(err)
		}

		ServiceMetricMap["/servicesname"+"/"+ns] = string(jsonServices)
	}
	return ServiceMetricMap

}
