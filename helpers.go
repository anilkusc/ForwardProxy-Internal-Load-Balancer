package main

import (
	"net"
	"net/http"
	"strings"
	"time"
)

func ClientCreator(ip string) *http.Transport {

	localAddr, err := net.ResolveIPAddr("ip", ip)
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
