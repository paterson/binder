package main

import (
    "github.com/gin-gonic/gin"
)

var store Store

func main() {
    r := gin.Default()
    router.POST("/signup", signup)
    router.POST("/login", login)
    router.Run(":3001")
    store := DefaultStore
}

// Take in username and password, and store in db as (username, encrypt(password))
func signup(ctx *gin.Context) {
    err := store.CreateUser(ctx.Query("username", ctx.Query("password"))
    if err != nil {
        token := GenerateToken()
        response := AuthenticatedResponse{token: token}
        encryptedResponse := response.encrypt(user.Password)
        ctx.JSON(code, encryptedResponse.EncodeJSON())
    } else {
        // Fail
    }
}

// Find row with username in db and ensure encrypt(password) = encrypted_password
func login(ctx *gin.Context) {
    exists := store.UserExists(ctx.Query("username", ctx.Query("password"))
    if exists {
        token := GenerateToken()
        response := AuthenticatedResponse{token: token}
        encryptedResponse := response.encrypt(user.Password)
        ctx.JSON(code, encryptedResponse.EncodeJSON())
    } else {
        // Fail
    }
}

func sampleMethod(ctx *gin.Context) {
    request, err := NewRequest(ctx)
    if err == nil {
        // Auth is valid
        filepath := request.Query("filepath")
        request.Respond(200, json)
    }
}
