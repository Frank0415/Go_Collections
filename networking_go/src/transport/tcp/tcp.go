package tcp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type Handler interface {
	ServeTCP(conn net.Conn, ctx context.Context, connID int32)
}

type Server struct {
	handler  Handler
	port     int
	maxConns int
}

func NewServer(handler Handler, port int, maxConns int) *Server {
	if handler == nil || port <= 0 || maxConns <= 0 || maxConns > 10000 || port > 65535 {
		log.Fatal("Invalid parameters to create TCP server.")
	}
	return &Server{
		handler:  handler,
		port:     port,
		maxConns: maxConns,
	}
} // After this we assume the server is valid

func bindPort(port int) (*net.TCPListener, error) {
	for currPort := port; currPort < port+100; currPort++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", currPort))
		if err != nil {
			return nil, fmt.Errorf("Bad TCP address formatting")
		}
		listener, err := net.ListenTCP("tcp", addr)
		if err == nil {
			log.Println("Now TCP listens on port", currPort)
			return listener, nil
		}
	}
	return nil, fmt.Errorf("could not bind to any port in range %d-%d", port, port+99)
}

func (s *Server) StartServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	listener, err := bindPort(s.port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	sem := make(chan struct{}, s.maxConns)
	var wg sync.WaitGroup
	var connID atomic.Int32
	connID.Store(0)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// 如果是关闭了 listener 导致的报错，说明是正常退出
				if errors.Is(err, net.ErrClosed) {
					return
				} else {
					log.Println("Accept error:", err)
					continue
				}
			}
			select {
			case sem <- struct{}{}:
				wg.Add(1)
				go func(c net.Conn) {
					defer func() { <-sem }()
					defer wg.Done()
					id := connID.Add(1)
					s.handler.ServeTCP(c, ctx, id)
				}(conn)
			default:
				conn.Write([]byte("Server busy, try again later.\n"))
				conn.Close()
			}
		}
	}()

	// the thread will wait here until sigChan receives a signal
	<-sigChan
	log.Println("Shutting down TCP server, waiting all connections to close.")
	listener.Close()
	cancel()  // cancel the context to notify all goroutines (if they are using it)
	wg.Wait() // the main goroutine will need to wait for all connection goroutines to finish
	log.Println("All connections closed, TCP server exited.")
}
