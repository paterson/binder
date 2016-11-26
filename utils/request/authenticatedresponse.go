package request

import (
    "github.com/gin-gonic/gin"
)

type AuthenticatedResponse struct {
    Token Token
}

type EncryptedAuthenticatedResponse struct {
    Token EncryptedToken
}

func (encryptedResponse  EncryptedAuthenticatedResponse) EncodeJSON() gin.H {
    token := encryptedResponse.Token
    return gin.H{
        "ticket":          token.Ticket.SessionKey,
        "session_key":     token.SessionKey,
        "server_identity": token.ServerIdentity,
        "timeout":         token.Timeout,
    }
}

func (encryptedResponse EncryptedAuthenticatedResponse) Decrypt(password string) AuthenticatedResponse {
    return AuthenticatedResponse{Token: encryptedResponse.Token.decrypt(password)}
}

func (response AuthenticatedResponse) Encrypt(password string) EncryptedAuthenticatedResponse {
    return EncryptedAuthenticatedResponse{Token: response.Token.encrypt(password)}
}
