package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var num int
	_, err := fmt.Scan(&num)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(num)
	if 10.0 < num {
		test, err := exec.Command("/bin/bash", "/home/shisui/DoanDuong/custom-scheduler/test_script.sh").Output()
		if err != nil {
			fmt.Println("Error:",err)
		}
		fmt.Printf(string(test))
	}
}
