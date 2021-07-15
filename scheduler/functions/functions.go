package functions

import (
	"encoding/json"
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

func GetHttpApi(domain string, query string, metrics MetricResponse) MetricResponse {
	resp, err := http.Get("http://" + domain + "api/v1/query?query=" + query + "")
	if err != nil {
		fmt.Println(err)
	}
	DecodeJsonDataToStruct(&metrics, resp)

	return metrics
}

func DecodeJsonDataToStruct(metrics *MetricResponse, resp *http.Response) {
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ConvertStringToFloat(metric MetricResponse) float64 {
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