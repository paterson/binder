package request

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

/* Request. Key: Session key */
type Request struct {
	Ticket Ticket
	Body   Body
	ctx    *gin.Context
}

type EncryptedRequest struct {
	Ticket EncryptedTicket
	body   EncryptedBody
}

func Authenticate(ctx *gin.Context) (Request, error) {
	ticket := EncryptedTicket{SessionKey: EncryptedSessionKey(ctx.Param("ticket"))}
	encryptedRequest := EncryptedRequest{Ticket: ticket, body: EncryptedBody(ctx.Param("body"))}
	request, err := encryptedRequest.decrypt()
	request.ctx = ctx
	return request, err
}

func (request Request) Method() string {
	if request.ctx.Request.Method == "POST" {
		return "POST"
	}
	return "GET"
}

func (request Request) Param(key string) string {
	if request.Method() == "GET" {
		return request.ctx.Query(key)
	} else {
		return request.ctx.PostForm(key)
	}
}

func (request Request) SendFile(filepath string) {
	request.ctx.File(filepath)
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
	ticket := encryptedRequest.Ticket.Decrypt()
	body, err := encryptedRequest.body.Decrypt(ticket.SessionKey)
	request := Request{Ticket: ticket, Body: body}
	return request, err
}

func (request Request) encrypt(sessionKey SessionKey) EncryptedRequest {
	return EncryptedRequest{Ticket: request.Ticket.Encrypt(), body: request.Body.Encrypt(sessionKey)}
}
