package http

import (
	"context"
	"net"
	"time"
)

const maxCacheSize = 1048576 // 1MB
const cacheBufferSize = 4096 // 4KB

type HTTPServer struct {
	// Add fields as necessary
}

func (httpserver HTTPServer) ServeTCP(conn net.Conn, ctx context.Context, connID int32) {
	defer conn.Close()

	buf := make([]byte, 4096)
	var reqBuf []byte

	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if len(reqBuf)+n > maxCacheSize {
			return
		}
		reqBuf = append(reqBuf, buf[:n]...)

		for {
			result := ParseInput(reqBuf)
			if result.Status == StatusIncomplete {
				break
			}

			if result.Status == StatusError {
				return
			}

			if result.Status == StatusComplete {
				ServeRequest(conn, result.Request.URI)
				LogInput(result.Request, true)

				// Keep-Alive 判断逻辑
				shouldClose := false

				// HTTP/1.0 默认关闭，除非明确声明 keep-alive
				if result.Request.Version == "HTTP/1.0" {
					if result.Request.Headers["Connection"] != "keep-alive" {
						shouldClose = true
					}
				} else {
					// HTTP/1.1 默认保持，但客户端可以要求关闭
					if result.Request.Headers["Connection"] == "close" {
						shouldClose = true
					}
				}

				if shouldClose {
					return // defer conn.Close() 会执行
				}

				// 继续处理剩余数据或等待下一个请求
				reqBuf = result.Remaining
				reqBuf = ShrinkBuffer(reqBuf, cacheBufferSize)

				// 如果还有剩余数据（粘包），立即处理下一个请求，不等待 Read
				if len(reqBuf) > 0 {
					continue
				}
			}
		}
	}
}
