package main

import (
	// "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func getHttpApi(domain string, query string) MetricResponse {
	resp, err := http.Get("http://"+domain+"api/v1/query?query="+query+"")
	if err != nil {
		fmt.Println(err)
	}
	var metrics MetricResponse
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

	var nodeMemMetric MetricResponse
	nodeMemMetric = getHttpApi("localhost:2505/", "((node_memory_MemTotal_bytes{job=\"node-exporter\"}-node_memory_MemAvailable_bytes{job=\"node-exporter\"})/(node_memory_MemTotal_bytes{job=\"node-exporter\"}))*100")

	var nodeCpuMetric MetricResponse
	nodeCpuMetric = getHttpApi("localhost:2505/", "100-irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\",cpu=\"1\"}[10m])*100")

	var nodeDiskMetric MetricResponse
	nodeDiskMetric = getHttpApi("localhost:2505/", "(node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node-exporter\"})/(1024*1024*1024)")

	var nodeReceiveNet MetricResponse
	nodeReceiveNet = getHttpApi("localhost:2505/", "irate(node_network_receive_bytes_total{device=\"eth0\"}[5m])/1024")

	var nodeTransmitNet MetricResponse
	nodeTransmitNet = getHttpApi("localhost:2505/", "irate(node_network_transmit_bytes_total{device=\"eth0\"}[5m])/1024")

	// Iterate through the metric results to find the node with the best value
	maxDisk := 0
	bestNode := ""
	for _, m := range nodeDiskMetric.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Disk size Value (MB): %s\n\n", m.MetricValue[1])

		// Convert string in metric results to an integer
		memDataConvert := fmt.Sprintf("%v", m.MetricValue[1])
		memData, err := strconv.ParseFloat(memDataConvert, 64)
		if err != nil {
			return "", err
		}

		switch m.MetricInfo["instance"] {
		case "10.47.0.10:9100":
			m.MetricInfo["instance"] = "node7"
		case "10.44.0.8:9100":
			m.MetricInfo["instance"] = "node6"
		}

		if memData < maxDisk {
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

