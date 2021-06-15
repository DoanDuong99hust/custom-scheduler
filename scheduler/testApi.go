package main

import (
	// "bufio"
	"encoding/json"
	// "errors"
	"fmt"
	"net/http"
	"os"
	// "strconv"
)

type MetricResponseTest struct {
	Data DataTest `json:"data,omitempty"`
}

type DataTest struct {
	Results []ResultTest `json:"result,omitempty"`
}

// Idea to use interface for metric values (which have different types) from
// https://stackoverflow.com/questions/38861295/how-to-parse-json-arrays-with-two-different-data-types-into-a-struct-in-go-lang
type ResultTest struct {
	MetricInfo  map[string]string `json:"metric,omitempty"`
	MetricValue []interface{}     `json:"value,omitempty"` //Index 0 is unix_time, index 1 is sample_value (metric value)
}

func decodeJsonDataToStructTest(metrics *MetricResponseTest, resp *http.Response) {
	decoder := json.NewDecoder(resp.Body)
	fmt.Println(decoder)
	err := decoder.Decode(metrics)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	respCpu, err3 := http.Get("http://localhost:8080/api/v1/query?query=100 - (avg by (instance) (irate(node_cpu_seconds_total{mode='idle'}[10m])) * 100)")
	if err3 != nil {
		fmt.Println(err3)
	}

	var metrics MetricResponseTest
	// defer respCpu.Body.Close()
	// err := json.NewDecoder(respCpu.Body).Decode(metrics)
	// fmt.Println(err)
	decodeJsonDataToStructTest(&metrics, respCpu)

	for _, m := range metrics.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Value: %s\n", m.MetricValue[1])
	}


}