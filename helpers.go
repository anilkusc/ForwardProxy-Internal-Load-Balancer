package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func ClientCreator(params ...string) (*http.Transport, string) {
	var nci string
	if len(params) <= 0 {
		nci = IpSelector()
	} else {
		nci = params[0]
	}

	localAddr, err := net.ResolveIPAddr("ip", nci)
	if err != nil {
		panic(err)
	}
	localTCPAddr := net.TCPAddr{
		IP: localAddr.IP,
	}

	d := net.Dialer{
		LocalAddr: &localTCPAddr,
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                d.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return tr, nci
}

func InitializeAddresses() []string {
	var addresses []string
	var address string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						address = addr.String()
						address = address[0:strings.Index(address, "/")]
						addresses = append(addresses, address)
					}
				}

			}
		}
	}
	addresses = RemoveDuplicateValues(addresses)

	return addresses
}

func IpSelector() string {
	var temp int
	rand.Seed(time.Now().UnixNano())
	if counter >= len(addresses) {
		counter = counter - len(addresses)
	}

	switch *LBAlghorithm {
	case "roundrobin":
		temp = counter
		counter++
		break
	case "random":
		temp = rand.Intn(len(addresses))
		break
	case "sourceip":
		temp = SourceIpHelper()
		break
	case "leastconn":
		temp = LeastConnHelper()
		break
	default:
		temp = counter
		counter++
		break
	}

	return addresses[temp]
}

func Healthcheck_Address() {
	for {
		if len(Bannedaddresses) > 0 {
			for i, Bannedaddress := range Bannedaddresses {

				tr, _ := ClientCreator(Bannedaddress)
				client := &http.Client{Transport: tr}
				req, err := http.NewRequest("GET", *checkAddr, nil)
				if err != nil {
					continue
				}
				resp, err := client.Do(req)
				if err != nil {
					continue
				}

				if resp.StatusCode < 400 {
					log.Println("This address is now healty and enabling again: ", Bannedaddress)
					addresses = append(addresses, Bannedaddress)
					Bannedaddresses = RemoveIndex(Bannedaddresses, i)
				}
			}
		}
		for i, address := range addresses {

			tr, _ := ClientCreator(address)
			client := &http.Client{Transport: tr}
			req, _ := http.NewRequest("GET", *checkAddr, nil)
			resp, err := client.Do(req)

			if err != nil || resp.StatusCode >= 400 {
				log.Println("Cannot reach to target via this interface:", address)
				Bannedaddresses = append(Bannedaddresses, address)
				addresses = RemoveIndex(addresses, i)
			}
		}
		if len(addresses) < 1 {
			log.Println("There is no available interface address!")
		}
		fmt.Println("Healty adresses are: ", addresses)
		time.Sleep(time.Duration(*timeInterval) * time.Second)

	}
}

//Remove index from a slice
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func RemoveDuplicateValues(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func AddBlacklist(ip string, host string) {
	inList := false
	for _, Bannedaddress := range Bannedaddresses {

		if Bannedaddress == ip {
			inList = true
			break
		}

	}
	if inList == false {
		Bannedaddresses = append(Bannedaddresses, ip)
	}
	for i, address := range addresses {
		if address == ip {
			addresses = RemoveIndex(addresses, i)
			break
		}

	}
	go Blacklist_Controller(ip, host)
}

func Blacklist_Controller(ip string, host string) {
	for {
		time.Sleep(300 * time.Second)
		tr, _ := ClientCreator(ip)
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode < 400 {
			log.Println("This address is now healty and enabling again: ", ip)
			addresses = append(addresses, ip)
			for i, Bannedaddress := range Bannedaddresses {
				if Bannedaddress == ip {
					Bannedaddresses = RemoveIndex(Bannedaddresses, i)
					break
				}
			}

		}

	}
}

func ExcludeAddresses(excluded string) []string {
	excludedAddresses := strings.Split(excluded, ",")
	newAddresses := addresses

	for _, excludedAddress := range excludedAddresses {
		for i, address := range addresses {

			if excludedAddress == address {
				newAddresses = RemoveIndex(newAddresses, i)
				break
			}

		}
	}
	return newAddresses
}
func LeastConnHelper() int {
	min := -1
	inAddresses, inValues := false, false
	var addressStore string
	for _, address := range addresses {
		inAddresses = false
		for key, _ := range LeastConnValues {
			if key == address {
				inAddresses = true
				break
			}

		}
		if inAddresses == false {
			LeastConnValues[address] = 0
		}
	}

	for key, _ := range LeastConnValues {
		inValues = false
		for _, address := range addresses {
			if key == address {
				inValues = true
				break
			}
		}
		if inValues == false {
			//Delete from leastconn values
			delete(LeastConnValues, key)
		}

	}
	for address, count := range LeastConnValues {
		if min == -1 {
			min = count
			addressStore = address
		} else {

			if count <= min {
				min = count
				addressStore = address
			}

		}
	}
	var index int
	for i, address := range addresses {
		if address == addressStore {
			index = i
			break
		}
	}
	LeastConnValues[addressStore]++
	return index
}
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	address := strings.Split(IPAddress, ":")
	if address[0] == "[" && address[2] == "1]" {
		address[0] = "127.0.0.1"
	}
	return address[0]
}

func SourceIpHelper() int {
	var isBanned bool = true

	if SourceIpCache[clientIP] == "" || SourceIpCache[clientIP] == "127.0.0.1" {
		temp := counter
		SourceIpCache[clientIP] = addresses[temp]
		counter++
		return temp
	} else {
		for _, address := range addresses {
			if address == SourceIpCache[clientIP] {
				isBanned = false
				break
			}
		}
	}

	if isBanned == true {
		temp := counter
		SourceIpCache[clientIP] = addresses[temp]
		counter++
		return temp
	} else {
		for i, address := range addresses {
			if address == SourceIpCache[clientIP] {
				return i
			}
		}

	}
	return 0
}
func ApiLogCollector(w *http.Response, r *http.Request, respBody string, reqBody string) {
	var request Request
	var response Response
	request.HttpVersion = r.Proto
	request.Host = r.Host + r.URL.Path
	request.Method = r.Method
	request.Body = string(reqBody)
	responseHeaderMap := make(map[string]string)
	for name, values := range r.Header {
		for _, value := range values {
			responseHeaderMap[name] = value
		}
	}
	request.Headers = responseHeaderMap
	response.HttpVersion = w.Proto
	response.Status = string(w.Status)
	response.Body = respBody
	requestHeaderMap := make(map[string]string)
	for name, values := range w.Header {
		for _, value := range values {
			requestHeaderMap[name] = value
		}
	}
	response.Headers = requestHeaderMap
	log := &Log{
		LogRequest:  request,
		LogResponse: response,
	}
	logString, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err)
	}
	if len(apiLogs) >= *logCacheCount {
		for len(apiLogs) >= *logCacheCount {
			apiLogs = RemoveIndex(apiLogs, 0)
		}
	}
	apiLogs = append(apiLogs, string(logString))

	f, _ := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(time.Now().UTC().String() + "\t" + string(logString) + "\n")
	fi, err := os.Stat("access.log")
	if err != nil {
		fmt.Println(err)
	}
	f.Close()
	size := fi.Size()
	if size/(1024*1024) > *logSize {
		err := os.Remove("access.log")
		if err != nil {
			fmt.Println(err)
		}
	}

}
