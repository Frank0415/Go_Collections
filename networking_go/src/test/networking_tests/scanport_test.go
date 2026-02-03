package main

import (
	"fmt"
	"net"
	"testing"
)

func TestScanLocalPortsAsync(t *testing.T) {
	// 启动一个监听器来测试扫描逻辑
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("无法启动监听器: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	fmt.Printf("测试正在运行，监听端口: %d\n", port)

	// 调用扫描函数
	// 目前 ScanLocalPortsAsync 会扫描 0-65535，这在测试中可能较慢
	// 但为了演示测试程序，我们直接调用它。
	// 建议实际开发中修改函数以支持指定端口范围。
	ScanLocalPortsAsync("tcp", "127.0.0.1", 100)
}
