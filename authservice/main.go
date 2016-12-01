package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/authservice/store"
	"github.com/paterson/binder/utils/logger"
	"github.com/paterson/binder/utils/request"
	"os"
)

var Store *store.Store

func main() {
	Store = store.DefaultStore()
	gin.DefaultWriter = logger.New("../log/authservice.log")
	router := gin.Default()
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.Run(port())
}

// Take in username and password, and store in db as (username, encrypt(password))
func signup(ctx *gin.Context) {
	user := store.User{Username: ctx.Param("username"), Password: ctx.Param("password")}
	Store = Store.CreateUser(&user)
	if Store.Result.Error == nil {
		token := request.GenerateToken()
		response := request.AuthenticatedResponse{Token: token}
		encryptedResponse := response.Encrypt(ctx.Param("password"))
		ctx.JSON(request.StatusOK, encryptedResponse.EncodeJSON())
	} else {
		ctx.JSON(request.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

// Find row with username in db and ensure encrypt(password) = encrypted_password
func login(ctx *gin.Context) {
	user := store.User{Username: ctx.Query("username"), Password: ctx.Query("password")}
	Store = Store.UserExists(&user)
	if Store.Result.SuccessfulQuery {
		token := request.GenerateToken()
		response := request.AuthenticatedResponse{Token: token}
		encryptedResponse := response.Encrypt(ctx.Query("password"))
		ctx.JSON(request.StatusOK, encryptedResponse.EncodeJSON())
	} else {
		ctx.JSON(request.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

func port() string {
	return ":" + os.Args[1]
}
