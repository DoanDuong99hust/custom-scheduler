package main

import (
	functions "customScheduler/scheduler/functions"
	"fmt"
	"log"
	"os/exec"
)

func deployment(service string, threshold float64) (float64, float64) {
	// lay thong so bang thong
	// neu bang thong k dat yeu cau se goi script
	var nodeReceiveNet functions.MetricResponse
	nodeReceiveNet = functions.GetHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
	bandwidth := functions.ConvertStringToFloat(nodeReceiveNet)

	// index := strings.Split("\"-\"", "-")
	// deploymentFile := index[0] + "./reDeploy" + service + ".sh" + index[0]
	switch service {
	case "Decode":
		_, err := exec.Command("/bin/bash", "./test_script.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		break
	case "Density":
		_, err := exec.Command("/bin/bash", "./test_script_1.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		break
	}



	return bandwidth, threshold

}

func reDeploy(service string)  {
	switch service {
	case "Decode":
		_, err := exec.Command("/bin/bash", "./updateDeployment.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		break
	case "Density":
		_, err := exec.Command("/bin/bash", "./updateDeployment1.sh").Output()
		if err != nil {
			log.Println("Error:", err)
		}
		break
	}
}
func main() {

	serviceName, serviceThreshold := functions.InputServiceReqired()
	fmt.Println(serviceName)
	fmt.Println(serviceThreshold)

	bandwidth, threshold := deployment(serviceName, serviceThreshold)

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
					reDeploy(serviceName)
					fmt.Println("-------------------")
				} else {
					continue
				}
			} else { // instantBw < threshold
				if preBandwidth >= threshold {
					fmt.Println("Previous Bandwidth: ", preBandwidth)
					fmt.Println("Instant Bandwidth: ", instantBw)
					reDeploy(serviceName)
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
					reDeploy(serviceName)
					fmt.Println("-------------------")
				} else {
					continue
				}
			} else {
				if preBandwidth < threshold {
					fmt.Println("Previous Bandwidth: ", preBandwidth)
					fmt.Println("Instant Bandwidth: ", instantBw)
					reDeploy(serviceName)
					fmt.Println("-------------------")
				} else {
					continue
				}
			}
		}
		preBandwidth = instantBw
	}
}
