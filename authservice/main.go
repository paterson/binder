package authservice

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/paterson/binder/authservice/store"
)

const (
    SERVER_KEY = "lkm4iuPKCCJQGBGB"
)

var Store *store.Store

func main() {
    Store = store.DefaultStore()
    router := gin.Default()
    router.POST("/signup", signup)
    router.POST("/login", login)
    router.Run(":3001")
}

// Take in username and password, and store in db as (username, encrypt(password))
func signup(ctx *gin.Context) {
    user := store.User{Username: ctx.Query("username"), Password: ctx.Query("password")}
    Store = Store.CreateUser(&user)
    if Store.Error == nil {
        token := GenerateToken()
        response := AuthenticatedResponse{token: token}
        encryptedResponse := response.encrypt(ctx.Query("password"))
        ctx.JSON(http.StatusOK, encryptedResponse.EncodeJSON())
    } else {
         ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
    }
}

// Find row with username in db and ensure encrypt(password) = encrypted_password
func login(ctx *gin.Context) {
    user := store.User{Username: ctx.Query("username"), Password: ctx.Query("password")}
    Store = Store.UserExists(&user)
    if Store.SuccessfulQuery {
        token := GenerateToken()
        response := AuthenticatedResponse{token: token}
        encryptedResponse := response.encrypt(ctx.Query("password"))
        ctx.JSON(http.StatusOK, encryptedResponse.EncodeJSON())
    } else {
         ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
    }
}
