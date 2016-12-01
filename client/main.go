package main

import (
	"github.com/paterson/binder/utils/clientproxy"
)

func main() {
	clientproxy := clientproxy.New()
	clientproxy.Login("niall", "password")
	clientproxy.WriteFile("test.png", "/test.png")
	clientproxy.ReadFile("/test.png", "test1.png")
}
