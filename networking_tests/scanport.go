package main

import (
	"fmt"
	"net"
	"time"
)

func scanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func ScanLocalPorts(){
	for i := range 65536{
		if scanPort("tcp","127.0.0.1",i) {
			fmt.Printf("Port %d is open.\n", i)
		}
	}
}