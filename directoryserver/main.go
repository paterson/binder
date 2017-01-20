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
	router.POST("/request/lock", lockRequest)
	router.POST("/request/unlock", unlockRequest)
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
		key := string(req.Ticket.SessionKey)
		Store = Store.EnsureHostExistsForPath(filepath)
		Store = Store.IsValidLockKeyForPath(key, filepath)
		if Store.Result.Error == nil && Store.Result.ValidLockKey {
			req.Respond(request.StatusOK, request.Params{"hosts": Store.Result.Hosts})
		} else {
			fmt.Println(request.Params{"error": "error"})
			req.Respond(request.Status400, request.Params{"error": "Something went wrong"})
		}
	}
}

func lockRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Params["filepath"]
		Store = Store.GetLockStatusForPath(filepath)
		if Store.Result.Locked {
			req.Respond(request.Status400, request.Params{"error": "File is already locked"})
		} else {
			key := string(req.Ticket.SessionKey)
			Store = Store.LockPathWithKey(filepath, key)
			req.Respond(request.StatusOK, request.Params{"success": "true"})
		}
	}
}

func unlockRequest(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil {
		filepath := req.Params["filepath"]
		key := string(req.Ticket.SessionKey)
		Store = Store.UnlockPathWithKey(filepath, key)
		if !Store.Result.Locked {
			req.Respond(request.StatusOK, request.Params{"success": "true"})
		} else {
			req.Respond(request.Status400, request.Params{"error": "Something went wrong"})
		}
	}
}

func port() string {
	return ":" + os.Getenv("PORT")
}
