package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/controller/.kube/config", "location to kubeconfig file")
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)
	pods, _ := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Println("Pod name: ", pod.Name)
		fmt.Println("Pod IP: ", pod.Status.PodIP)
		fmt.Println("Host IP", pod.Status.HostIP)
	}
}
