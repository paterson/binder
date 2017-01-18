package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/authservice/store"
	"github.com/paterson/binder/utils/request"
	"os"
)

var Store *store.Store

func main() {
	Store = store.DefaultStore()
	router := gin.Default()
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.Run(port())
}

// Take in username and password, and store in db as (username, encrypt(password))
func signup(ctx *gin.Context) {
	user := store.User{Username: ctx.PostForm("username"), Password: ctx.PostForm("password")}
	Store = Store.CreateUser(&user)
	if Store.Result.Error == nil {
		token := request.GenerateToken()
		response := request.AuthenticatedResponse{Token: token}
		encryptedResponse := response.Encrypt(user.Password)
		ctx.JSON(request.StatusOK, encryptedResponse.EncodeJSON())
	} else {
		ctx.JSON(request.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

// Find row with username in db and ensure encrypt(password) = encrypted_password
func login(ctx *gin.Context) {
	user := store.User{Username: ctx.PostForm("username"), Password: ctx.PostForm("password")}
	Store = Store.UserExists(&user)
	if Store.Result.SuccessfulQuery {
		fmt.Println("Successful Query")
		token := request.GenerateToken()
		response := request.AuthenticatedResponse{Token: token}
		encryptedResponse := response.Encrypt(user.Password)
		ctx.JSON(request.StatusOK, encryptedResponse.EncodeJSON())
	} else {
		fmt.Println("Not Successful Query")
		ctx.JSON(request.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

func port() string {
	return ":" + os.Getenv("PORT")
}
