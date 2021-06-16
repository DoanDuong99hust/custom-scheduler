package main

import (
	// "bufio"
	"encoding/json"
	"io/ioutil"
	// "reflect"

	// "log"

	// "errors"
	"fmt"
	"net/http"
	"os"

	// "golang.org/x/text/number"
	"strconv"
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

// type MetricValue struct {
// 	MetricValue		float64 		
// 	Status 			string			
// }

func decodeJsonDataToStructTest(metrics *MetricResponseTest, resp *http.Response) {
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// NewReader := bufio.NewReader(os.Stdin)
	// fmt.Print("Insert Prometheus Domain: ")
	// promDomain,_ := NewReader.ReadString('/')
	// respCpu, err3 := http.Get("http://localhost:8080/api/v1/query?query=100-irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\",cpu=\"1\"}[10m])*100")
	// if err3 != nil {
	// 	fmt.Println(err3)
	// }

	respMem, err2 := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable_bytes{job=\"node-exporter\"}/(1024*1024)")
	if err2 != nil {
		fmt.Println(err2)
	}
	
	var metrics MetricResponseTest
	decodeJsonDataToStructTest(&metrics, respMem)
	
	for _, m := range metrics.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Value: %s\n", m.MetricValue[1])
		memData := fmt.Sprintf("%v",m.MetricValue[1])
		metricValue, err := strconv.ParseFloat(memData,64)
		if err != nil {
			fmt.Println(err) 
		}
		fmt.Println(metricValue-100)
	}
	
	file,_:= json.MarshalIndent(metrics," ", " ")
	_ = ioutil.WriteFile("nodeCpu.json", file, 0644)
}