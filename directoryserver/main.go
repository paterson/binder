package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/directoryserver/store"
	"github.com/paterson/binder/utils/request"
	"os"
)

var Store *store.Store

func main() {
	Store = store.DefaultStore()
	Store.Seed()
	router := gin.Default()
	router.POST("/request/read", readRequest)
	router.POST("/request/write", writeRequest)
	router.Run(port())
}

func readRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Params["filepath"]
		Store = Store.HostsForPath(filepath)
		if Store.Result.Error == nil {
			req.Respond(request.StatusOK, request.Params{"hosts": Store.Result.Hosts})
		} else {
			req.Respond(request.Status404, request.Params{"error": "404"})
		}
	}
}

func writeRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Params["filepath"]
		Store = Store.EnsureHostExistsForPath(filepath)
		if Store.Result.Error == nil {
			req.Respond(request.StatusOK, request.Params{"hosts": Store.Result.Hosts})
		} else {
			fmt.Println(request.Params{"error": "error"})
			req.Respond(request.Status400, request.Params{"error": "Something went wrong"})
		}
	}
}

func port() string {
	return ":" + os.Args[1]
}
