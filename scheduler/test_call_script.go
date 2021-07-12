package main

// import (
// 	"fmt"
// 	"os/exec"
// 	"strconv"
// )

// func main() {

// 	var net MetricResponse
// 	net = getHttpApi("localhost:8080/", "irate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[5m])/(1024*1024)", net)

// 	for _, m := range net.Data.Results {

// 		memDataConvert := fmt.Sprintf("%v", m.MetricValue[1])
// 		bandwidth, err := strconv.ParseFloat(memDataConvert, 64)
	
// 		if  bandwidth > 10 {
			
// 			test, err := exec.Command("/bin/bash", "./test_script.sh").Output()
// 			if err != nil {
// 				fmt.Println("Error:",err)
// 			}
// 			fmt.Printf(string(test))
// 		}
// 	}
// }