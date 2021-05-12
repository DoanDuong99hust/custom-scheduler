package main

import (
	"errors"
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

func getSecondBestNode(nodes []Node) (string, error){
	var machineStatus MachineStatus
	worker01 := getMongoDbData(machineStatus, "worker01")
	worker02 := getMongoDbData(machineStatus, "worker02")

	nodeStatus := []MachineStatus{worker01,worker02}
	min := worker01

	for _, node := range nodeStatus {
		nodeMemory := node.Cpu[1]
		if (min.FreeDisk[1]/(1024*1024)) > 400 {
			if nodeMemory <= min.Cpu[1] {
				min = node
			}
		}
	}
	// machineStatus = getMongoDbData(machineStatus, "shisui")
	bestNode := min.Machine

	if bestNode == "" {
		return "", errors.New("No node found")
	} else {
		return bestNode, nil
	}
}

// func main() {
// 	var nodes []Node
// 	var bestNodeName, _ = getSecondBestNode(nodes)
// 	fmt.Println(bestNodeName)
// }