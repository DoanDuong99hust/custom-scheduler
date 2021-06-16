package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Struct for decoded JSON from HTTP response
type MetricResponse struct {
	Data Data `json:"data,omitempty"`
}

type Data struct {
	Results []Result `json:"result,omitempty"`
}

// Idea to use interface for metric values (which have different types) from
// https://stackoverflow.com/questions/38861295/how-to-parse-json-arrays-with-two-different-data-types-into-a-struct-in-go-lang
type Result struct {
	MetricInfo  map[string]string `json:"metric,omitempty"`
	MetricValue []interface{}     `json:"value,omitempty"` //Index 0 is unix_time, index 1 is sample_value (metric value)
}

type Layer1 struct {
	Status 		string		`json:"status,omitempty"`
	Data		Data		`json:"data,omitempty"`
}
// Returns the name of the node with the best metric value
func getBestNodeName(nodes []Node) (string, error) {
	var nodeNames []string
	for _, n := range nodes {
		nodeNames = append(nodeNames, n.Metadata.Name)
	}

	// Add prometheus domain
	fmt.Print("Insert Prometheus Domain: ")
	proDomain, _ := bufio.NewReader(os.Stdin).ReadString('/')

	// Execute a query over the HTTP API to get the metric kubelet_node_name
	respName, err := http.Get("http://" + proDomain + "/api/v1/query?query=kubelet_node_name")
	if err != nil {
		fmt.Println(err)
	}

	// Decode the JSON body of the HTTP response into a struct
	var nodeNameMetric MetricResponse
	decodeJsonDataToStruct(&nodeNameMetric, respName)

	// Execute a query over the HTTP API to get the metric node_memory
	respMem, err2 := http.Get("http://" + proDomain + "/api/v1/query?query=((avg_over_time(node_memory_MemTotal_bytes{job=\"node-exporter\"}[5m]) - avg_over_time(node_memory_MemFree_bytes[5m]) - avg_over_time(node_memory_Cached_bytes[5m]) - avg_over_time(node_memory_Buffers_bytes[5m])) / avg_over_time(node_memory_MemTotal_bytes[5m])) * 100")
	if err2 != nil {
		fmt.Println(err2)
	}

	// Decode the JSON body of the HTTP response into a struct
	var nodeMemMetric MetricResponse
	decodeJsonDataToStruct(&nodeMemMetric, respMem)

	// Execute a query over the HTTP API to get the metric node_cpu_seconds_total
	respCpu, err3 := http.Get("http://"+proDomain+"/api/v1/query?query=100-irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\",cpu=\"1\"}[10m])*100")
	if err3 != nil {
		fmt.Println(err3)
	}

	// Decode the JSON body of the HTTP response into a struct
	var nodeCpuMetric MetricResponse
	decodeJsonDataToStruct(&nodeCpuMetric, respCpu)

	respDisk, err4 := http.Get("http://"+proDomain+"/api/v1/query?query=(node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node-exporter\"})/(1024*1024)")
	if err4 != nil {
		fmt.Println(err4)
	}

	// Decode the JSON body of the HTTP response into a struct
	var nodeDiskMetric MetricResponse
	decodeJsonDataToStruct(&nodeDiskMetric, respDisk)

	// Iterate through the metric results to find the node with the best value
	max := 0
	bestNode := ""
	for _, m := range nodeDiskMetric.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Value: %s\n\n", m.MetricValue[1])

		// Convert string in metric results to an integer
		metricValue, err := strconv.Atoi(m.MetricValue[1].(string))
		if err != nil {
			return "", err
		}

		if metricValue > max {
			// Check if the node is in the list passed in (nodes the pod will fit on)
			available := nodeAvailable(nodeNames, m.MetricInfo["instance"])
			if available == true {
				max = metricValue
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
	fmt.Println(decoder)
	err := decoder.Decode(metrics)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


// func main() {
	// fmt.Print("Insert Prometheus Domain: ")
	// proDomain, _:= bufio.NewReader(os.Stdin).ReadString('/')

	// resp, err1 := http.Get("http://localhost:8080/api/v1/query?query=kubelet_node_name")
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }

	// respMem, err2 := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable_bytes{job=\"node-exporter\"}/(1024*1024)")
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }

	// respCpu, err3 := http.Get("http://localhost:8080/api/v1/query?query=100 - (avg by (instance, job) (irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\"}[10m])) * 100)")
	// if err3 != nil {
	// 	fmt.Println(err3)
	// }
	
	// respDisk, err4 := http.Get("http://localhost:8080/api/v1/query?query=node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node-exporter\"}")
	// if err4 != nil {
	// 	fmt.Println(err4)
	// }

	// var nodeNameMetric MetricResponse
	// decodeJsonDataToStruct(&nodeNameMetric, resp)

	// for _, m := range nodeNameMetric.Data.Results{
	// 	fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
	// 	fmt.Printf("Value: %s\n", m.MetricValue[1])
	// }

	// var metrics MetricResponse
	// decodeJsonDataToStruct(&metrics, respCpu)

	// for _, m := range metrics.Data.Results {
	// 	// Print metric value for the node
	// 	fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
	// 	fmt.Printf("Value: %s\n", m.MetricValue[1])
	// }

// }
