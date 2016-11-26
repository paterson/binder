package authservice

import (
    "github.com/gin-gonic/gin"
)

type AuthenticatedResponse struct {
    token Token
}

type EncryptedAuthenticatedResponse struct {
    token EncryptedToken
}

func (encryptedResponse  EncryptedAuthenticatedResponse) EncodeJSON() gin.H {
    token := encryptedResponse.token
    return gin.H{
        "ticket":          token.Ticket.SessionKey,
        "session_key":     token.SessionKey,
        "server_identity": token.ServerIdentity,
        "timeout":         token.Timeout,
    }
}

func (encryptedResponse EncryptedAuthenticatedResponse) decrypt(password string) AuthenticatedResponse {
    return AuthenticatedResponse{token: encryptedResponse.token.decrypt(password)}
}

func (response AuthenticatedResponse) encrypt(password string) EncryptedAuthenticatedResponse {
    return EncryptedAuthenticatedResponse{token: response.token.encrypt(password)}
}
