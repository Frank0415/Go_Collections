package http

var mimeTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".css":  "text/css",
	".js":   "application/javascript",
	".jpg":  "image/jpeg", // 注意：jpeg 的 MIME 是 image/jpeg，不是 image/jpg
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".ico":  "image/x-icon",
}

type Filetype int

const (
	html Filetype = iota
	css
	js
	jpg
	jpeg
	png
	gif
	ico
	unknown
)

var FileMimeTypes = map[Filetype]string{
	html: "text/html; charset=utf-8",
	css:  "text/css",
	js:   "application/javascript",
	jpg:  "image/jpeg", // 注意：jpeg 的 MIME 是 image/jpeg，不是 image/jpg
	jpeg: "image/jpeg",
	png:  "image/png",
	gif:  "image/gif",
	ico:  "image/x-icon",
	unknown: "application/octet-stream",
}