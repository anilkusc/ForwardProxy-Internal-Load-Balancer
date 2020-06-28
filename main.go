package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addresses = InitializeAddresses()
	counter   = 0
	port      = flag.String("port", "8080", "Specify port number")
)

func proxy(w http.ResponseWriter, req *http.Request) {
	counter++
	if counter >= len(addresses) {
		counter = counter - len(addresses)
	}

	tr := ClientCreator(addresses[counter])
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(req.Method, "http://"+req.Host, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error reading response. ", err)

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	w.Write(body)
}

func main() {
	flag.Parse()
	//fmt.Println(addresses)

	http.HandleFunc("/", proxy)

	http.ListenAndServe(":"+*port, nil)
}
