package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addresses       = InitializeAddresses()
	Bannedaddresses = []string{}
	counter         = 0
	port            = flag.String("port", "8080", "Specify port number")
	checkAddr       = flag.String("check-addr", "", "Healthcheck for specific address if address is not reachable from interface don't use that network card for that specific address.")
	checkInterface  = flag.Bool("check-interface", false, "If some of the interfaces are down they are disabled for proxy until they get healty")
	timeInterval    = flag.Int("time-interval", 300, "Healthcheck time interval in seconds.")
)

func proxy(w http.ResponseWriter, req *http.Request) {

	tr := ClientCreator()
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
	if *checkAddr != "" {
		go Healthcheck_Address()
	}
	//fmt.Println(addresses)

	http.HandleFunc("/", proxy)
	fmt.Println("Serving on port :", *port)
	http.ListenAndServe(":"+*port, nil)
}
