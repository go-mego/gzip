# Gzip

Gzip 套件可以透過 Gzip 壓縮網頁內容，這能夠有效減少網路流量，但有可能稍微提高 CPU 用量（取決於時機）。

# 索引

* [安裝方式](#安裝方式)
* [使用方式](#使用方式)
    * [單一路由](#單一路由)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/go-mego/gzip
```

# 使用方式

將 `gzip.Zipper` 傳入 `Use` 來在所有路由中啟用 Gzip 壓縮。

```go
package main

import (
	"github.com/go-mego/gzip"
	"github.com/go-mego/mego"
)

func main() {
    m := mego.New()
    // 將 Gzip 壓縮技術套用到全域路由上。
	m.Use(gzip.Zipper())
	m.Run()
}
```

## 單一路由

Gzip 中介軟體也能夠套用到單一路由而非所有路由。

```go
// 將 Gzip 壓縮技術套用到指定的單個路由。
m.Get("/", gzip.Zipper(), func() string {
    return "哈囉，世界！"
})
```