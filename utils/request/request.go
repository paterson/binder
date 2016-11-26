package request

import (
    "github.com/gin-gonic/gin"
    "mime/multipart"
)

/* Request. Key: Session key */
type Request struct {
    Ticket Ticket
    Body   Body
    ctx    *gin.Context
}

type EncryptedRequest struct {
    Ticket          EncryptedTicket
    body            EncryptedBody
}

func Authenticate(ctx *gin.Context) (Request, error) {
    ticket := EncryptedTicket{SessionKey: EncryptedSessionKey(ctx.Query("ticket"))}
    encryptedRequest := EncryptedRequest{Ticket: ticket, body: EncryptedBody(ctx.Query("body"))}
    request, err := encryptedRequest.decrypt()
    request.ctx = ctx
    return request, err
}

func (request Request) Query(key string) string {
    return request.ctx.Query(key)
}

func (request Request) RetrieveUploadedFile() (multipart.File, string, error) {
    file, header, err := request.ctx.Request.FormFile("upload")
    filename := header.Filename
    return file, filename, err
}

func (request Request) Respond(code int, body Body) {
    response := Response{Body: body}
    encryptedResponse := response.encrypt(request.Ticket.SessionKey)
    request.ctx.JSON(code, encryptedResponse.EncodeJSON())
}

func (encryptedRequest EncryptedRequest) decrypt() (Request, error) {
    // decrypt request by decrypting ticket with server key, then decrypt body with session key in ticket
    ticket    := encryptedRequest.Ticket.Decrypt()
    body, err := encryptedRequest.body.Decrypt(ticket.SessionKey)
    request   := Request{Ticket: ticket, Body: body}
    return request, err
}

func (request Request) encrypt(sessionKey SessionKey) EncryptedRequest {
    return EncryptedRequest{Ticket: request.Ticket.Encrypt(), body: request.Body.Encrypt(sessionKey)}
}
