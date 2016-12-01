package request

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

/* Request. Key: Session key */
type Request struct {
	Ticket Ticket
	Params Params
	ctx    *gin.Context
}

type EncryptedRequest struct {
	Ticket EncryptedTicket
	Params EncryptedParams
}

func Authenticate(ctx *gin.Context) (Request, error) {
	ticket := EncryptedTicket{SessionKey: EncryptedSessionKey(ctx.Param("ticket"))}
	encryptedRequest := EncryptedRequest{Ticket: ticket, Params: NewEncryptedParams(ctx)}
	request, err := encryptedRequest.decrypt()
	request.ctx = ctx
	return request, err
}

func (request Request) SendFile(filepath string) {
	request.ctx.File(filepath)
}

func (request Request) RetrieveUploadedFile() (multipart.File, string, error) {
	file, header, err := request.ctx.Request.FormFile("file")
	filename := header.Filename
	return file, filename, err
}

func (request Request) Respond(code int, params Params) {
	response := Response{Params: params}
	encryptedResponse := response.encrypt(request.Ticket.SessionKey)
	request.ctx.JSON(code, encryptedResponse.EncodeJSON())
}

func (encryptedRequest EncryptedRequest) decrypt() (Request, error) {
	ticket := encryptedRequest.Ticket.Decrypt()
	params, err := encryptedRequest.Params.Decrypt(ticket.SessionKey)
	request := Request{Ticket: ticket, Params: params}
	return request, err
}

func (request Request) encrypt(sessionKey SessionKey) EncryptedRequest {
	return EncryptedRequest{Ticket: request.Ticket.Encrypt(), Params: request.Params.Encrypt(sessionKey)}
}
