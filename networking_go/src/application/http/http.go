package http

import (
	"context"
	"net"
)

type HTTPServer struct {
	// Add fields as necessary
}

func (httpserver *HTTPServer) ServeTCP(conn net.Conn, ctx context.Context, connID int32) {

}

func StartHTTPServer() {}
