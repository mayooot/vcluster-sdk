package main

import (
	"context"
	"github.com/mayooot/vcluster-sdk/pkg/connection"

	"github.com/loft-sh/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	// using in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building clientset: %s", err.Error())
	}

	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		klog.Fatalf("Error getting server version: %s", err.Error())
	} else {
		klog.Info("Successfully connected to Root Kubernetes cluster")
	}

	vClient, err := connection.GetVClusterClientset(context.TODO(), clientset, "test", "vcluster-test", log.GetInstance())
	if err != nil {
		klog.Fatalf("Error building vcluster clientset: %s", err.Error())
	}

	podList, err := vClient.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		klog.Fatalf("Error listing pods: %s", err.Error())
	}
	for _, pod := range podList.Items {
		klog.Infof("Pod name: %s, namespace: %s", pod.Name, pod.Namespace)
	}

	select {}
}
