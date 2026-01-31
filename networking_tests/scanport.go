package main

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
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

func ScanLocalPorts() {
	for i := range 65536 {
		if scanPort("tcp", "127.0.0.1", i) {
			fmt.Printf("Port %d is open.\n", i)
		}
	}
}

func ScanLocalPortsAsync(protocol, hostname string, workerCount int) {
	var wg sync.WaitGroup
	const rangeIP = 65536

	if workerCount > runtime.NumCPU()*10 {
		workerCount = runtime.NumCPU() * 10
	}

	semaphore := make(chan struct{}, workerCount)
	results := make(chan int, rangeIP)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for port := range rangeIP {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			dialer := net.Dialer{}
			conn, err := dialer.DialContext(ctx, protocol, net.JoinHostPort(hostname, fmt.Sprintf("%d", port)))
			if err == nil {
				conn.Close()
				results <- port
			}
		}(port)
	}

	wg.Wait()
	close(results)

	for port := range results {
		fmt.Printf("Port %d is open.\n", port)
	}
}

func GetAllInterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range interfaces {
		fmt.Printf("接口: %s, MAC: %s\n", iface.Name, iface.HardwareAddr)

		// 获取接口地址
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			fmt.Printf("  地址: %s\n", addr.String())
		}
	}
}
