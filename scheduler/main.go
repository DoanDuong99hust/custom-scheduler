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

func main() {
	log.Printf("Starting %s scheduler...", schedulerName)

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	log.Println("b1")

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)

	log.Println("b2")

	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)

	log.Println("b3")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("b4")

	for {
		log.Println("loop")
		// lay thong so bang thong
		// neu bang thong k dat yeu cau se goi script
		var nodeReceiveNet MetricResponse
		nodeReceiveNet = getHttpApi("localhost:8080/", "irate(node_network_receive_bytes_total{device=\"enp0s3\",instance=\"192.168.101.192:9100\"}[5m])/(1024*1024)", nodeReceiveNet)
		bandwidth := convertStringToFloat(nodeReceiveNet)
		log.Println(bandwidth)
		if bandwidth < 1.15 {
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
		
		time.Sleep(2)
		// khi chuong trinh bi dung lai
		select {	
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}
