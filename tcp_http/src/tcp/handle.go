package tcp

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	var mainBuf []byte
	buf := make([]byte, 1024)

	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Read timeout, closing")
				return
			}
			if err == io.EOF {
				log.Println("Connection closed by client")
				if len(mainBuf) > 0 {
					log.Printf("Connection Closed with remaining %q", mainBuf)
				}
			} else {
				log.Println("Read error:", err)
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
