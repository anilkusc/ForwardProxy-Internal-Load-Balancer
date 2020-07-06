package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addresses          = InitializeAddresses()
	clientIP           = ""
	Bannedaddresses    = []string{}
	counter            = 0
	LeastConnValues    = make(map[string]int)
	SourceIpCache      = make(map[string]string)
	port               = flag.String("port", "8080", "Specify port number")
	checkAddr          = flag.String("check-addr", "", "Healthcheck for specific address if address is not reachable from interface don't use that network card for that specific address.")
	LBAlghorithm       = flag.String("balancing-alghorithm", "roundrobin", "Specify Round-Robin Alghorithm.(roundrobin,random,sourceip,leastconn)")
	checkInterface     = flag.Bool("check-interface", false, "If some of the interfaces are down they are disabled for proxy until they get healty")
	timeInterval       = flag.Int("time-interval", 300, "Healthcheck time interval in seconds.")
	excludedeaddresses = flag.String("excluded-addresses", "", "Specify ip addresses that exclude for load balancing.E.g. :`192.168.1.20,192.168.1.21` ")
)

func proxy(w http.ResponseWriter, req *http.Request) {
	if *LBAlghorithm == "sourceip" {
		clientIP = ReadUserIP(req)
	}
	for _, _ = range addresses {
		tr, ip := ClientCreator()
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest(req.Method, "http://"+req.Host, nil)
		if err != nil {
			fmt.Println("This ip address is now in blacklist:", ip)
			AddBlacklist(ip, "http://"+req.Host)
			log.Println(err)
			counter = 0
			continue
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Println("Error doing request. ", err)
			fmt.Println("This ip address is now in blacklist:", ip)
			AddBlacklist(ip, "http://"+req.Host)
			counter = 0
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response. ", err)
			fmt.Println("This ip address is now in blacklist:", ip)
			AddBlacklist(ip, "http://"+req.Host)
			counter = 0
			continue
		}
		defer resp.Body.Close()
		w.Write(body)
		break
	}

}

func main() {
	flag.Parse()
	addresses = ExcludeAddresses(*excludedeaddresses)
	if *checkAddr != "" {
		go Healthcheck_Address()
	}

	http.HandleFunc("/", proxy)
	fmt.Println("Serving on port :", *port)
	http.ListenAndServe(":"+*port, nil)
}
