package main

import (
    "github.com/gin-gonic/gin"
)

/* Response. Key: Session key */
type Response struct {
    body  Body
}

type EncryptedResponse struct {
    body EncryptedBody
}

func (encryptedResponse EncryptedResponse) EncodeJSON() gin.H {
    return gin.H{
        "body": encryptedResponse.body
    }
}

func (encryptedResponse EncryptedResponse) decrypt(sessionKey string) Response {
    // TODO
}

func (response Response) encrypt(sessionKey string) EncryptedResponse {
    // encrypt body with SessionKey as key.
}
