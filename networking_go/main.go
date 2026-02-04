package main

import "tcp_http/src/test/tcp_server"
import "tcp_http/src/application/http"
import "tcp_http/src/transport/tcp"

func main() {
	// orig_test()
	http_test()
}

func orig_test() {
	tcp_server.StartTCPServer()
}

func http_test() {
	var httpserver http.HTTPServer
	server := tcp.NewServer(httpserver, 10000, 100)
	server.StartServer()
}