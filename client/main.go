package main

import (
	"github.com/paterson/binder/clientproxy"
)

func main() {
	clientproxy := clientproxy.New()
	clientproxy.WriteFile("test.png", "/test.png")
	clientproxy.ReadFile("/test.png", "test1.png")
}
