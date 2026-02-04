package http

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

func fetchBody(path string, code int) (io.ReadCloser, int64, error) {
	if code != 200 {
		msg := fmt.Sprintf("<h1>%d %s</h1>", code, getStatusText(code))
		return io.NopCloser(strings.NewReader(msg)), int64(len(msg)), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, err
	}

	return file, info.Size(), nil
}

func identifyFiletype(ext string) Filetype {
	switch ext {
	case ".html":
		return html
	case ".css":
		return css
	case ".js":
		return js
	case ".jpg", ".jpeg":
		return jpg
	case ".png":
		return png
	case ".gif":
		return gif
	case ".ico":
		return ico
	default:
		return unknown
	}
}

func ServeRequest(w io.Writer, uri string) {
	filePath, fileType, statusCode := fetchResources(uri)

	bodyReader, size, err := fetchBody(filePath, statusCode)
	if err != nil {
		statusCode = 500
		bodyReader, size, _ = fetchBody("", statusCode)
	}
	defer bodyReader.Close()

	contentType := FileMimeTypes[fileType]
	if statusCode != 200 {
		contentType = "text/html; charset=utf-8"
	}

	// 写入状态行
	w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, getStatusText(statusCode))))
	// 写入响应头
	w.Write([]byte(fmt.Sprintf("Content-Type: %s\r\n", contentType)))
	w.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", size)))
	w.Write([]byte("\r\n"))

	// 使用 io.Copy 将 body 写入连接，适合流式传输（图片、大文件等）
	io.Copy(w, bodyReader)
}

func fetchResources(uri string) (string, Filetype, int) {
	absDocumentRoot, err := filepath.Abs("resources")
	if err != nil {
		return "", unknown, 500
	}

	// 安全第一步：路径清理
	path := filepath.Clean(filepath.Join(absDocumentRoot, uri))

	// 第二步：路径安全（生死线）
	// 攻击向量：GET ../../etc/passwd HTTP/1.1
	// 防御逻辑：验证边界，确保 path 处于 absDocumentRoot 之内
	if !strings.HasPrefix(path, absDocumentRoot) {
		return "", unknown, 403
	}

	return checkFile(path)
}

func checkFile(path string) (string, Filetype, int) {
	info, err := os.Stat(path)
	if err != nil {
		// 错误映射
		if os.IsNotExist(err) {
			return "", unknown, 404
		}
		if os.IsPermission(err) {
			return "", unknown, 403
		}
		// 其他 I/O 错误
		return "", unknown, 500
	}

	// 拒绝目录（或根据逻辑尝试 index.html）
	if info.IsDir() {
		indexPath := filepath.Join(path, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return indexPath, html, 200
		}
		// 禁止直接访问目录
		return "", unknown, 403
	}

	return path, identifyFiletype(filepath.Ext(path)), 200
}
