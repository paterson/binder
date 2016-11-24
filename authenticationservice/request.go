package main

import (
    "github.com/gin-gonic/gin"
)

type Body            map[string]string
type EncryptedBody   string

/* Request. Key: Session key */
type Request struct {
    ticket EncryptedTicket
    Body   Body
    ctx    *gin.Context
}

type EncryptedRequest struct {
    ticket          EncryptedTicket
    body            EncryptedBody
}

func NewRequest(ctx *gin.Context) (Request, error) {
    ticket := EncryptedTicket{SessionKey: ctx.Query("ticket")}
    encryptedRequest := EncryptedRequest{ticket: ticket, body: ctx.Query("body")} // TODO structs not strings..
    request, err := encryptedRequest.decrypt()
    request.ctx = ctx
    return request, err
}

func (request Request) Query(key string) string {
    return request.ctx.Query(key)
}

func (request Request) Respond(code int, body Body) {
    response := Response{body: body}
    encryptedResponse := response.encrypt(request.ticket.SessionKey)
    request.ctx.JSON(code, encryptedResponse.EncodeJSON())
}

func (encryptedRequest EncryptedRequest) decrypt() (Request, error) {
    // decrypt request by decrypting ticket with server key, then decrypt body with session key in ticket
    ticket    := encryptedRequest.ticket.decrypt()
    body, err := encryptedRequest.body.decrypt(ticket.SessionKey) // pseudo
    request   := Request{ticket: encryptedRequest.ticket, body: body}
    return request, err
}

func (request Request) encrypt() EncryptedRequest {
    // TODO
}
