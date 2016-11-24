package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    router.POST("/signup", signup)
    router.POST("/login", login)
}

// Take in username and password, and store in db as (username, encrypt(password))
func signup(ctx *gin.Context) {
    // Create user from params (add sessionKey to db?)
    // If successful, create request with user and Context and send response
    token := GenerateToken()
    response := AuthenticatedResponse{token: token}
    encryptedResponse := response.encrypt(user.Password)
    ctx.JSON(code, encryptedResponse.EncodeJSON())
}

// Find row with username in db and ensure encrypt(password) = encrypted_password
func login(ctx *gin.Context) {
    // Find user from params
    // If successful, create request with user and Context and send response
    token := GenerateToken()
    response := AuthenticatedResponse{token: token}
    encryptedResponse := response.encrypt(user.Password)
    ctx.JSON(code, encryptedResponse.EncodeJSON())
}

func sampleMethod(ctx *gin.Context) {
    request, err := NewRequest(ctx)
    if err == nil {
        // Auth is valid
        filepath := request.Query("filepath")
        request.Respond(200, json)
    }
}
