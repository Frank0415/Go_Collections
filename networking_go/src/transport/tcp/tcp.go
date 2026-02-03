package tcp

import (
	"context"
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
	handler  *Handler
	port     int
	maxConns int
}

func NewServer(handler *Handler, port int, maxConns int) *Server {
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

func StartServer(server *Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	listener, err := bindPort(server.port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	sem := make(chan struct{}, server.maxConns)
	var wg sync.WaitGroup
	var connID atomic.Int32
	connID.Store(0)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// 如果是关闭了 listener 导致的报错，说明是正常退出
				if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
					return
				} else {
					log.Println("Accept error:", err)
					continue
				}
			}
			select {
			case sem <- struct{}{}:
				wg.Add(1)
				go func() {
					defer func() { <-sem }()
					defer wg.Done()
					connID.Add(1)
					(*server.handler).ServeTCP(conn, ctx, connID.Load())
				}()
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

/*
0.1版本：
	设计的时候的一个想法是用
		done := make(chan struct{})

		go func() {
			<-sigChan
			log.Println("Shutting down TCP server, waiting all connections to close.")
			close(done)
			listener.Close()
		}()
	然后再在后面的函数里面接受这个done然后再去用这个done来控制退出，
	比如说handleTCPConnection里面也接受这个done然后退出

	但是由于一开始做的初始化的问题，还是先做一个简单的版本，
	直接在这个循环里面用这个sigchan，也免得需要把主循环套在goroutine里面

0.2版本：
	for {
	log.Println("Head of 'for' loop")
	select {
	case <-sigChan:
		log.Println("Shutting down TCP server, waiting all connections to close.")
		listener.Close()
		wg.Wait()
		log.Println("All connections closed, TCP server exited.")
		return
	default:
		conn, err := listener.Accept()
		if err != nil {
			// 如果是关闭了 listener 导致的报错，说明是正常退出
			select {
			case <-sigChan:
				// wg.Wait()
				// log.Println("All connections closed, TCP server exited.")
				return
			default:
				log.Println("Accept error:", err)
				continue
			}
		}
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func() {
				defer func() { <-sem }()
				defer wg.Done()
				handleTCPConnection(conn)
			}()
		default:
			conn.Write([]byte("Server busy, try again later.\n"))
			conn.Close()
		}
	}
	这个版本的一个重大问题， select 结构看似在检查信号，但逻辑是这样的：
		循环开始，select 检查 sigChan —— 没信号
		进入 default 分支，执行 listener.Accept()
		Accept() 阻塞了！ 程序卡在这里等待新连接，无法回到 select 重新检查信号
		你按 Ctrl+C，信号被操作系统发送给 sigChan，但代码还在第 3 步卡着，无法执行
		直到有新连接进来（或出错），Accept() 返回，代码才回到循环开头，发现 sigChan 有信号
	所以这个版本的代码实际上是无法优雅退出的，学到了一切阻塞的东西都要放到 goroutine 里面去做这个教训

0.3版本：
	还是需要goroutine,部分回滚到0.1,但是尝试自己设计，不依赖LLM
*/
