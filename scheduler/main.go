// Original file obtained from https://github.com/kelseyhightower/scheduler/blob/master/main.go
// Modified to change scheduler name as indicated in inline comments below
// -------------------------------------------------------------------------------------------
// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -------------------------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
	
)

// Modified scheduler name
const schedulerName = "customScheduler"

func deployment(service string, threshold float64 ) (float64, float64) {
	// lay thong so bang thong
	// neu bang thong k dat yeu cau se goi script
	var nodeReceiveNet MetricResponse
	nodeReceiveNet = getHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
	bandwidth := convertStringToFloat(nodeReceiveNet)

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

func main() {
	log.Printf("Starting %s scheduler...", schedulerName)

	doneChan := make(chan struct{})
	var wg sync.WaitGroup
	
	fmt.Println("b1")
	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)

	fmt.Println("b2")
	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)

	fmt.Println("b3")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("b4")
	
	bandwidth,threshold := deployment("Decode", 0.87)
	fmt.Println(bandwidth)
	fmt.Println("b5")
	
	if bandwidth < threshold { 
		preInstant := threshold
		for {
			var nodeReceiveNet MetricResponse
			nodeReceiveNet = getHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
			instantBw := convertStringToFloat(nodeReceiveNet)

			//fmt.Println("Previous Bandwidth: ", bandwidth)
			//fmt.Println("Instant Bandwidth: ", instantBw)
			//fmt.Println("------------------")
			if instantBw >=threshold {
				if instantBw > preInstant {
				        fmt.Println("Previous Bandwidth: ", preInstant)
				        fmt.Println("Instant Bandwidth: ", instantBw)
				     
					deployment("Decode", threshold)
					fmt.Println("-------------------")
				} else {
					continue
				}	
			} else if instantBw < threshold {
				if instantBw < preInstant {
				        fmt.Println("Previous Bandwidth: ", preInstant)
				        fmt.Println("Instant Bandwidth: ", instantBw)
				        
					deployment("Decode", threshold)
					fmt.Println("-------------------")
				} else {
					continue
				}
			} 
			preInstant = instantBw
			time.Sleep(30)

		}
	} else {
		preInstant := threshold
		for {
			var nodeReceiveNet MetricResponse
			nodeReceiveNet = getHttpApi("localhost:8080/", "rate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[1m])/(1024*1024)", nodeReceiveNet)
			instantBw := convertStringToFloat(nodeReceiveNet)
			if instantBw < preInstant {
		                fmt.Println("Previous Bandwidth: ", bandwidth)
		                fmt.Println("Instant Bandwidth: ", instantBw)
		                fmt.Println("---------")
				deployment("Decode", threshold)
			} else {
				//fmt.Println("continue...")
				continue
			}
			preInstant = instantBw
			time.Sleep(30)

		}
	}

	// for {
	// 	// khi chuong trinh bi dung lai
	// 	select {
	// 	case <-signalChan:
	// 		log.Printf("Shutdown signal received, exiting...")
	// 		close(doneChan)
	// 		wg.Wait()
	// 		os.Exit(0)
	// 	}
	// }
}
