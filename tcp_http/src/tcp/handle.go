package tcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func handleTCPConnection(conn net.Conn, ctx context.Context, connID int32) {
	defer conn.Close()

	var mainBuf []byte
	buf := make([]byte, 1024)

	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, err := conn.Read(buf)

		if err != nil {
			var opErr *net.OpError
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-ctx.Done():
					log.Println("Context cancelled, closing connection, connID:", connID)
					conn.Write([]byte("connection closing, connID: " + fmt.Sprint(connID) + "\n"))
				default:
					continue
				}
			} else if errors.As(err, &opErr) {
				log.Println("Received stopping signal, closing connID:", connID)

			} else if err == io.EOF {
				log.Println("Connection closed by client")

			} else {
				log.Printf("Read error: %T", err)
			}

			if len(mainBuf) > 0 {
				log.Printf("Connection Closed with remaining %q", mainBuf)
			}
			return
		}

		log.Printf("Received data: %s", string(buf[:n]))
		if len(mainBuf)+n > 1024*1024 { // 1MB
			log.Println("Message too long, closing connection")
			return
		}
		mainBuf = append(mainBuf, buf[:n]...)

		for {
			index := bytes.IndexByte(mainBuf, '\n')
			if index == -1 {
				break
			}
			line := mainBuf[:index+1]
			data := bytes.TrimRight(line, "\r\n")
			conn.Write([]byte("Echo: " + string(data) + "\n"))

			remaining := len(mainBuf) - (index + 1)
			if remaining > 0 {
				copy(mainBuf, mainBuf[index+1:])
				mainBuf = mainBuf[:remaining]
			} else {
				mainBuf = mainBuf[:0]
			}

			if cap(mainBuf) > len(mainBuf)*4 && cap(mainBuf) > 1024 {
				tmp := make([]byte, len(mainBuf))
				copy(tmp, mainBuf)
				mainBuf = tmp
			}

		}
		if len(mainBuf) > 1024*1024 { // 1MB
			log.Println("Message too long, closing connection")
			return
		}
	}
}

/*
原始设计：
使用atomic bool来标记当前是否在读数据，
在context取消的时候，等待读数据完成之后再关闭连接
关闭时开销较大，而且实现比较复杂，改为使用read deadline的方式来实现
read deadline的每秒loop一次的开销还需要测试
func handleTCPConnection(conn net.Conn, ctx context.Context) {
	defer conn.Close()

	var mainBuf []byte
	var stuckRead atomic.Bool
	buf := make([]byte, 1024)

	go func() {
		<-ctx.Done()
		for stuckRead.Load() != true {}
		conn.Close()
		log.Println("Context cancelled, closing connection")
	}()

	for {
		stuckRead.Store(true)

		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.Read(buf)

		stuckRead.Store(false)
		if err != nil {
			var opErr *net.OpError
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Read timeout, closing")

			} else if errors.As(err, &opErr) {
				log.Println("Received stopping signal, closing")

			} else if err == io.EOF {
				log.Println("Connection closed by client")

			} else {
				log.Printf("Read error: %T", err)
			}

			if len(mainBuf) > 0 {
				log.Printf("Connection Closed with remaining %q", mainBuf)
			}
			return
		}

		log.Printf("Received data: %s", string(buf[:n]))
		time.Sleep(10 * time.Second)
		if len(mainBuf)+n > 1024*1024 { // 1MB
			log.Println("Message too long, closing connection")
			return
		}
		mainBuf = append(mainBuf, buf[:n]...)

		for {
			index := bytes.IndexByte(mainBuf, '\n')
			if index == -1 {
				break
			}
			line := mainBuf[:index+1]
			data := bytes.TrimRight(line, "\r\n")
			conn.Write([]byte("Echo: " + string(data) + "\n"))

			remaining := len(mainBuf) - (index + 1)
			if remaining > 0 {
				copy(mainBuf, mainBuf[index+1:])
				mainBuf = mainBuf[:remaining]
			} else {
				mainBuf = mainBuf[:0]
			}

			if cap(mainBuf) > len(mainBuf)*4 && cap(mainBuf) > 1024 {
				tmp := make([]byte, len(mainBuf))
				copy(tmp, mainBuf)
				mainBuf = tmp
			}

		}
		if len(mainBuf) > 1024*1024 { // 1MB
			log.Println("Message too long, closing connection")
			return
		}
	}
}

*/
