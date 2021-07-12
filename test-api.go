// package main

// import (
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// func  main()  {
// 	resp, err := http.Get("http://127.0.0.1:5000//api/v1/resources/node-cpu")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//We Read the response body on the line below.
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//Convert the body to type string
// 	sb := string(body)
// 	log.Printf(sb)
// }
