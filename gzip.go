package gzip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-mego/mego"
)

type CompressionType int

const (
	CompressionTypeBest    CompressionType = gzip.BestCompression
	CompressionTypeFastest CompressionType = gzip.BestSpeed
	CompressionTypeDefault CompressionType = gzip.DefaultCompression
	CompressionTypeNone    CompressionType = gzip.NoCompression
)

type Options struct {
	Level CompressionType
}

func New(option ...*Options) mego.HandlerFunc {
	o := &Options{
		Level: CompressionTypeDefault,
	}
	if len(option) > 0 {
		o = option[0]
	}

	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, int(o.Level))
		if err != nil {
			panic(err)
		}
		return gz
	}
	return func(c *mego.Context) {
		if !shouldCompress(c.Request) {
			return
		}

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(c.Writer)

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{c.Writer, gz}
		defer func() {
			gz.Close()
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()
		c.Next()
	}
}

type gzipWriter struct {
	mego.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") {

		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
