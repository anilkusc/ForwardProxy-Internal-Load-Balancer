package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
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
		//
		break
	case "leastconn":
		//
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
