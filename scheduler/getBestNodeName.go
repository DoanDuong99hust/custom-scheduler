package main

import (
	// "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	// "os/exec"
)

func getHttpApi(domain string, query string, metrics MetricResponse) MetricResponse {
	resp, err := http.Get("http://" + domain + "api/v1/query?query=" + query + "")
	if err != nil {
		fmt.Println(err)
	}
	decodeJsonDataToStruct(&metrics, resp)

	return metrics
}

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

	var nodeDiskMetric MetricResponse
	nodeDiskMetric = getHttpApi("localhost:8080/", "(node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node_exporter_metrics\"})/(1024*1024*1024)", nodeDiskMetric)

	var nodeReceiveNet MetricResponse
	nodeReceiveNet = getHttpApi("localhost:8080/", "irate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[5m])/(1024*1024)", nodeReceiveNet)

	// var nodeTransmitNet MetricResponse
	// nodeTransmitNet = getHttpApi("localhost:2505/", "irate(node_network_transmit_bytes_total{device=\"eth0\"}[5m])/1024", nodeTransmitNet)

	// Iterate through the metric results to find the node with the best value
	maxDisk := 0.0
	bestNode := ""
	bandwidth := convertStringToFloat(nodeReceiveNet)
	if bandwidth >= 1.15 {
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
func decodeJsonDataToStruct(metrics *MetricResponse, resp *http.Response) {
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func convertStringToFloat(metric MetricResponse) float64 {
	for _, result := range metric.Data.Results {
		rawData := fmt.Sprintf("%v", result.MetricValue[1])
		convertedData, err := strconv.ParseFloat(rawData, 64)
		if err != nil {
			fmt.Println(err)
		}
		return convertedData
	}
	return 0.0
}
