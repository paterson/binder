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

func (encryptedResponse EncryptedAuthenticatedResponse) EncodeJSON() gin.H {
	return encryptedResponse.Token.ToJSON()
}

func (encryptedResponse EncryptedAuthenticatedResponse) Decrypt(password string) AuthenticatedResponse {
	return AuthenticatedResponse{Token: encryptedResponse.Token.Decrypt(password)}
}

func (response AuthenticatedResponse) Encrypt(password string) EncryptedAuthenticatedResponse {
	return EncryptedAuthenticatedResponse{Token: response.Token.Encrypt(password)}
}
