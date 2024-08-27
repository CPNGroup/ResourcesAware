package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func GetClientSet() (*kubernetes.Clientset, error) {
	kubeconfig := "./admin.conf"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	config.Host = "https://10.129.82.112:6443" // 替换为虚拟机IP
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create clientset: %v", err)
	}
	return clientset, nil
}

func GetConfig() (*rest.Config, error) {
	kubeconfig := "./admin.conf"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	config.Host = "https://10.129.82.112:6443" // 替换为虚拟机IP
	return config, nil
}
