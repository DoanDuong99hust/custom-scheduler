package main

import (
	functions "customScheduler/scheduler/functions"
	"fmt"
	"log"
	"os/exec"
)

const THRESHOLD = 0.87
func deployment(service string, threshold float64) (float64, float64) {
	// lay thong so bang thong
	// neu bang thong k dat yeu cau se goi script
	var nodeReceiveNet functions.MetricResponse
	nodeReceiveNet = functions.GetHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
	bandwidth := functions.ConvertStringToFloat(nodeReceiveNet)

	deploymentFile := "./reDeploy" + service + ".sh"
	fmt.Println(deploymentFile)
	if bandwidth < threshold {
		test, err := exec.Command("/bin/bash", "./test_script.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		log.Printf(string(test))
	} else {
		test, err := exec.Command("/bin/bash", "./test_script.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		log.Printf(string(test))
	}

	return bandwidth, threshold

}
func main()  {
	bandwidth, threshold := deployment("Decode", THRESHOLD)
	fmt.Println(bandwidth)

	preBandwidth := bandwidth
	for {
		var nodeReceiveNet functions.MetricResponse
		nodeReceiveNet = functions.GetHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
		instantBw := functions.ConvertStringToFloat(nodeReceiveNet)

		if preBandwidth < threshold {
			if instantBw >= threshold {
				if preBandwidth < threshold {
					fmt.Println("Previous Bandwidth: ", preBandwidth)
					fmt.Println("Instant Bandwidth: ", instantBw)
					deployment("Decode", THRESHOLD)
					fmt.Println("-------------------")
				} else {
					continue
				}
			} else { // instantBw < threshold
				if preBandwidth >= threshold {
					fmt.Println("Previous Bandwidth: ", preBandwidth)
					fmt.Println("Instant Bandwidth: ", instantBw)
					deployment("Decode", THRESHOLD)
					fmt.Println("-------------------")
				} else {
					continue
				}
			}
		} else { // bandwidth >= threshold
			if instantBw < threshold {
				if preBandwidth >= threshold {
					fmt.Println("Previous Bandwidth: ", preBandwidth)
					fmt.Println("Instant Bandwidth: ", instantBw)
					deployment("Decode", THRESHOLD)
					fmt.Println("-------------------")
				} else {
					if preBandwidth < threshold {
						fmt.Println("Previous Bandwidth: ", preBandwidth)
						fmt.Println("Instant Bandwidth: ", instantBw)
						deployment("Decode", THRESHOLD)
						fmt.Println("-------------------")
					} else {
						continue
					}
				}
			}
		}
		preBandwidth = instantBw
	}
}