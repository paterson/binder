package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/directoryservice/store"
)

var Store *store.Store

func main() {
	Store = store.DefaultStore()
	Store.CreateDefaultFileServerRecord()
	router := gin.Default()
	router.POST("/request/read", readRequest)
	router.POST("/request/write", writeRequest)
	router.Run(port())
}

func readRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Param("filepath")
		Store = Store.HostForPath(filepath)
		if Store.Result.Error == nil {
			host := Store.Result.Host
			request.Respond(request.StatusOK, request.Body{"host": host})
		} else {
			request.Respond(request.Status404, request.Body{"error": "404"})
		}
	}
}

func writeRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Param("filepath")
		Store = Store.EnsureHostExistsForPath(filepath)
		if Store.Result.Error == nil {
			host := Store.Result.Host
			request.Respond(request.StatusOK, request.Body{"host": host})
		} else {
			request.Respond(request.Status400, request.Body{"error": "Something went wrong"})
		}
	}
}

func port() string {
	return ":" + os.Args[1]
}
