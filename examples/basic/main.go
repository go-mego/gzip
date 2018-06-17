package main

import (
	"net/http"

	"github.com/go-mego/gzip"
	"github.com/go-mego/mego"
)

func main() {
	e := mego.Default()
	e.GET("/", gzip.New(), func(c *mego.Context) {
		c.String(http.StatusOK, "The string has been gzipped")
	})
	e.Run()
}
