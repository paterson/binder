package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/directoryservice/store"
	"github.com/paterson/binder/utils/logger"
	"github.com/paterson/binder/utils/request"
)

var Store *store.Store

func main() {
	Store = store.DefaultStore()
	Store.CreateDefaultFileServerRecord()
	gin.DefaultWriter = logger.DirectoryServiceLogger
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
			req.Respond(request.StatusOK, request.Body{"host": host})
		} else {
			req.Respond(request.Status404, request.Body{"error": "404"})
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
			req.Respond(request.StatusOK, request.Body{"host": host})
		} else {
			req.Respond(request.Status400, request.Body{"error": "Something went wrong"})
		}
	}
}

func port() string {
	return ":" + os.Args[1]
}
