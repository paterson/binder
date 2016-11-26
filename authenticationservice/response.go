package main

import (
    "github.com/gin-gonic/gin"
)

/* Response. Key: Session key */
type Response struct {
    Body  Body
}

type EncryptedResponse struct {
    body EncryptedBody
}

func (encryptedResponse EncryptedResponse) EncodeJSON() gin.H {
    return gin.H{
        "body": encryptedResponse.body,
    }
}

func (encryptedResponse EncryptedResponse) decrypt(sessionKey SessionKey) (Response, error) {
    body, err := encryptedResponse.body.Decrypt(sessionKey)
    return Response{Body: body}, err
}

func (response Response) encrypt(sessionKey SessionKey) EncryptedResponse {
    return EncryptedResponse{body: response.Body.Encrypt(sessionKey)}
}
