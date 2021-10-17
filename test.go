package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// "log"
	"net/http"
	// "net/url"
	// "bytes"
)

type Data struct {
	Switchs []Switch `json:"switch"`
}

type Switch struct {
	Ports []Port `json:"port"`
}
type Port struct {
	Tx float64 `json:"Tx"`
	Rx float64 `json:"Rx"`
}

func main() {
	resp, err := http.Get("http://localhost:5000/api/v1/test")
	if err != nil {
		fmt.Println(err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var metrics Data
	json.Unmarshal(responseData, &metrics)

	fmt.Println(metrics)
	for _, v := range metrics.Switchs {
		for _, j := range v.Ports {
			fmt.Println(j.Rx)
		}
	}

}
