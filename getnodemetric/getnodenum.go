package getnodemetric

import (
	"context"
	"pro/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func GetNodeResource() int {
	// 创建 Kubernetes 客户端
	clientset, err := common.GetClientSet()
	if err != nil {
		panic(err)
	}

	// 获取节点列表
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 获取节点数量
	nodeNum := len(nodes.Items)

	return nodeNum

}
