package main

import (
	"context"
	"encoding/json"
	"log"
	"pro/getnodemetric"
	"pro/getpodmetric"
	"pro/getservicename"
	"pro/listallnamespace"
	"pro/pb"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	address = "10.129.82.112:30030" // zookeeper地址
)

func WriteIn() {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewZkServiceClient(conn)

	nodeNum := getnodemetric.GetNodeResource()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.Set(ctx, &pb.PathAndData{Path: "/nodenum", Data: strconv.Itoa(nodeNum)})
	if err != nil {
		log.Printf("PodMetricMap写入错误: %v", err)
	}

	namespaceList := listallnamespace.ListAllNamespce()
	jsonNamespaceList, err := json.Marshal(namespaceList)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.Set(ctx, &pb.PathAndData{Path: "/namespaces", Data: string(jsonNamespaceList)})
	if err != nil {
		log.Printf("PodMetricMap写入错误: %v", err)
	}

	PodMetricMap := getpodmetric.GetPodMetric()
	for key, value := range PodMetricMap {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = c.Set(ctx, &pb.PathAndData{Path: key, Data: value})
		if err != nil {
			log.Printf("PodMetricMap写入错误: %v", err)
		}

	}

	PodResourceMap := getpodmetric.GetPodResources()
	for key, value := range PodResourceMap {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = c.Set(ctx, &pb.PathAndData{Path: key, Data: value})
		if err != nil {
			log.Printf("PodResourceMap写入错误: %v", err)
		}
	}

	ServiceMetricMap := getservicename.ListServices()
	for key, value := range ServiceMetricMap {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.Set(ctx, &pb.PathAndData{Path: key, Data: value})
		if err != nil {
			log.Printf("ServiceMetricMap写入错误: %v", err)
		}
	}

}

func main() {
	// 每隔30秒执行一次
	for {
		WriteIn()

		// 等待30秒
		time.Sleep(30 * time.Second)
	}

}
