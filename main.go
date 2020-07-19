package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	addresses          = InitializeAddresses()
	clientIP           = ""
	Bannedaddresses    = []string{}
	counter            = 0
	LeastConnValues    = make(map[string]int)
	SourceIpCache      = make(map[string]string)
	apiLogs            = []string{}
	port               = flag.String("port", "8080", "Specify port number")
	checkAddr          = flag.String("check-addr", "", "Healthcheck for specific address if address is not reachable from interface don't use that network card for that specific address.")
	LBAlghorithm       = flag.String("balancing-alghorithm", "roundrobin", "Specify Round-Robin Alghorithm.(roundrobin,random,sourceip,leastconn)")
	checkInterface     = flag.Bool("check-interface", false, "If some of the interfaces are down they are disabled for proxy until they get healty")
	timeInterval       = flag.Int("time-interval", 300, "Healthcheck time interval in seconds.")
	excludedeaddresses = flag.String("excluded-addresses", "", "Specify ip addresses that exclude for load balancing.E.g. :`192.168.1.20,192.168.1.21` ")
	logCacheCount      = flag.Int("cache-count", 50, "Specify log cache number of connections.")
	logSize            = flag.Int64("log-size", 50, "Specify log size(mb) to write log file.After 2. log file is reach this size previous log file will be removed")
)

func proxy(w http.ResponseWriter, req *http.Request) {
	if *LBAlghorithm == "sourceip" {
		clientIP = ReadUserIP(req)
	}
	for _, _ = range addresses {
		tr, ip := ClientCreator()
		client := &http.Client{Transport: tr}
		var reqBody string
		req, err := http.NewRequest(req.Method, "http://"+req.Host+req.URL.Path, nil)
		if err != nil {
			fmt.Println("This ip address is now in blacklist:", ip)
			AddBlacklist(ip, "http://"+req.Host)
			log.Println(err)
			counter = 0
			continue
		}
		if req.Body != nil {
			tempBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Println(err)
			}
			reqBody = string(tempBody)
		} else {
			reqBody = ""
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
		ApiLogCollector(resp, req, string(body), string(reqBody))
		w.Write(body)
		break
	}

}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
func api(w http.ResponseWriter, r *http.Request) {
	var connLogs string = `{ "Logs":[`
	for i, apiLog := range apiLogs {
		if i == len(apiLogs)-1 {
			connLogs = connLogs + apiLog
		} else {
			connLogs = connLogs + apiLog + ","
		}

	}
	connLogs = connLogs + `]}`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	fmt.Fprintf(w, connLogs)
}

func main() {
	flag.Parse()
	addresses = ExcludeAddresses(*excludedeaddresses)
	if *checkAddr != "" {
		go Healthcheck_Address()
	}
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsEndpoint)
	r.HandleFunc("/api", api)
	//If none of the above matches use as default below
	r.PathPrefix("/").HandlerFunc(proxy)
	fmt.Println("Serving on port :", *port)

	http.ListenAndServe(":"+*port, r)
}
