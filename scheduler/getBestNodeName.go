package main

import (
	functions "customScheduler/scheduler/functions"
	"strings"
	// "bufio"
	//"encoding/json"
	"errors"
	"fmt"

	//"net/http"
	//"os"
	"strconv"
	// "os/exec"
)

// Returns the name of the node with the best metric value
func getBestNodeName(nodes []Node) (string, error) {
	var nodeNames []string
	for _, n := range nodes {
		nodeNames = append(nodeNames, n.Metadata.Name)
	}

	// Add prometheus domain
	// fmt.Print("Insert Prometheus Domain: ")
	// proDomain, _ := bufio.NewReader(os.Stdin).ReadString('/')

	// var nodeMemMetric MetricResponse
	// nodeMemMetric = getHttpApi("localhost:4040/", "((node_memory_MemTotal_bytes{job=\"node-exporter\"}-node_memory_MemAvailable_bytes{job=\"node-exporter\"})/(node_memory_MemTotal_bytes{job=\"node-exporter\"}))*100", nodeMemMetric)

	// var nodeCpuMetric MetricResponse
	// nodeCpuMetric = getHttpApi("localhost:4040/", "100-irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\",cpu=\"1\"}[10m])*100", nodeCpuMetric)

	var nodeDiskMetric functions.MetricResponse
	nodeDiskMetric = functions.GetHttpApi("localhost:8080/", "(node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node_exporter_metrics\"})/(1024*1024*1024)", nodeDiskMetric)

	var nodeReceiveNet functions.MetricResponse
	nodeReceiveNet = functions.GetHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)

	// var nodeTransmitNet MetricResponse
	// nodeTransmitNet = getHttpApi("localhost:2505/", "irate(node_network_transmit_bytes_total{device=\"eth0\"}[5m])/1024", nodeTransmitNet)

	// Iterate through the metric results to find the node with the best value
	maxDisk := 0.0
	bestNode := ""
	var serviceArray []string
	bandwidth := functions.ConvertStringToFloat(nodeReceiveNet)
	serviceThreshold := 0.0
	_,_,podList := getUnscheduledPods()
	for _, pod := range podList.Items {
		index := strings.Split(pod.Metadata.Name, "-")
		podServiceName := index[0]
		if podServiceName == "decode" || 
		podServiceName == "density" || 
		podServiceName == "nginx" ||
		podServiceName == "nginx1" {			
			serviceArray = append(serviceArray, podServiceName)
		}
	}
	forLoop:for _, podServiceName := range serviceArray {
		fmt.Println("Pod name: ",podServiceName)
		switch podServiceName {
		case "nginx":
			fmt.Println(1)
			serviceThreshold = functions.DECODE_THRESHOLD
			break forLoop
		case "nginx1":
			fmt.Println(2)
			serviceThreshold = functions.DENSITY_THRESHOLD
			break forLoop
		}
	}
	fmt.Println("Bandwidth: ", bandwidth)
	fmt.Println("Threshold: ", serviceThreshold)
	if bandwidth >= serviceThreshold {
		for _, m := range nodeDiskMetric.Data.Results {
			// Print metric value for the node
			fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
			fmt.Printf("Disk size Value (GB): %s\n\n", m.MetricValue[1])

			// Convert string in metric results to an integer
			memDataConvert := fmt.Sprintf("%v", m.MetricValue[1])
			memData, err := strconv.ParseFloat(memDataConvert, 64)
			if err != nil {
				return "", err
			}

			switch m.MetricInfo["instance"] {
			case "192.168.101.191:9100":
				m.MetricInfo["instance"] = "edge"
			case "192.168.101.192:9100":
				m.MetricInfo["instance"] = "server"
			}

			if memData > maxDisk {
				// Check if the node is in the list passed in (nodes the pod will fit on)
				available := nodeAvailable(nodeNames, m.MetricInfo["instance"])
				fmt.Println(m.MetricInfo["instance"])
				fmt.Println("available is ", available)
				if available == true {
					maxDisk = memData
					bestNode = m.MetricInfo["instance"]
				}
			}
		}
	} else {
		bestNode = "edge"
	}

	if bestNode == "" {
		return "", errors.New("No node found")
	} else {
		return bestNode, nil
	}
}
func nodeAvailable(nodeNames []string, name string) (result bool) {
	for _, n := range nodeNames {
		fmt.Println(n + ", ")
		if name == n {
			return true
		}
	}
	return false
}

// Decode JSON data into a struct to get the metric values

