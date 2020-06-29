package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func ClientCreator() *http.Transport {

	localAddr, err := net.ResolveIPAddr("ip", IpSelector())
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
	return tr
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
	return addresses
}

func IpSelector() string {
	counter++
	if counter >= len(addresses) {
		counter = counter - len(addresses)
	}
	return addresses[counter]
}

func Healthcheck_Address() {
	for {
		if len(Bannedaddresses) > 0 {
			for i, Bannedaddress := range Bannedaddresses {

				localAddr, _ := net.ResolveIPAddr("ip", Bannedaddress)
				localTCPAddr := net.TCPAddr{
					IP: localAddr.IP,
				}

				d := net.Dialer{
					LocalAddr: &localTCPAddr,
				}
				tr := &http.Transport{
					Dial: d.Dial,
				}
				client := &http.Client{Transport: tr}
				req, _ := http.NewRequest("GET", *checkAddr, nil)
				resp, _ := client.Do(req)

				if resp.StatusCode < 400 {
					log.Println("This address is now healty and enabling again: ", Bannedaddress)
					addresses = append(addresses, Bannedaddress)
					RemoveIndex(Bannedaddresses, i)
				}
			}
		}
		for i, address := range addresses {
			localAddr, _ := net.ResolveIPAddr("ip", address)
			localTCPAddr := net.TCPAddr{
				IP: localAddr.IP,
			}

			d := net.Dialer{
				LocalAddr: &localTCPAddr,
			}
			tr := &http.Transport{
				Dial: d.Dial,
			}
			client := &http.Client{Transport: tr}
			req, _ := http.NewRequest("GET", *checkAddr, nil)
			resp, err := client.Do(req)

			if err != nil || resp.StatusCode >= 400 {
				log.Println("Cannot reach to target via this interface:", address)
				Bannedaddresses = append(Bannedaddresses, address)
				RemoveIndex(addresses, i)
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
