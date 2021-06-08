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
	Data Data `json:"data"`
}

type Data struct {
	Results []Result `json:"result"`
}

// Idea to use interface for metric values (which have different types) from
// https://stackoverflow.com/questions/38861295/how-to-parse-json-arrays-with-two-different-data-types-into-a-struct-in-go-lang
type Result struct {
	MetricInfo map[string]string  `json:"metric"`
	MetricValue []interface{} `json:"value"` //Index 0 is unix_time, index 1 is sample_value (metric value)
}

// Returns the name of the node with the best metric value
func getBestNodeName(nodes []Node) (string, error) {
	var nodeNames []string
	for _, n := range nodes {
		nodeNames = append(nodeNames, n.Metadata.Name)
	}

	// Add prometheus domain
	proDomain, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	// Execute a query over the HTTP API to get the metric node_memory_MemAvailable
	respName, err := http.Get("http://" + proDomain +"/api/v1/query?query=kubelet_node_name")
	if err != nil {
		fmt.Println(err)
	}
	// Decode the JSON body of the HTTP response into a struct
	var metrics MetricResponse
	decodeJsonDataToStruct(&metrics, respName)


	// Iterate through the metric results to find the node with the best value
	max := 0
	bestNode := ""
	for _, m := range metrics.Data.Results {
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
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {

	//resp, err1 := http.Get("http://localhost:8080/api/v1/query?query=kubelet_node_name")
	//if err1 != nil {
	//	fmt.Println(err1)
	//}

	respMem, err2 := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable_bytes{job=\"node-exporter\"}/(1024*1024)")
	if err2 != nil {
		fmt.Println(err2)
	}

	//respCpu, err3 := http.Get("http://localhost:8080/api/v1/query?query=(sum (rate (container_cpu_usage_seconds_total{image!=\"\"}[1m])) by (instance))*100")
	//if err3 != nil {
	//	fmt.Println(err3)
	//}
	//
	//respDisk, err4 := http.Get("http://localhost:8080/api/v1/query?query=node_filesystem_avail_bytes{mountpoint=\"/\",fstype=\"xfs\",job=\"node-exporter\"}")
	//if err4 != nil {
	//	fmt.Println(err4)
	//}
	var metrics MetricResponse
	decodeJsonDataToStruct(&metrics, respMem)

	for _, m := range metrics.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Value: %s\n", m.MetricValue[1])
	}

}
