package request

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var StatusOK = http.StatusOK
var StatusUnauthorized = http.StatusUnauthorized
var Status404 = 404
var Status400 = 400

/* Response. Key: Session key */
type Response struct {
	Params Params
}

type EncryptedResponse struct {
	Params EncryptedParams
}

func (encryptedResponse EncryptedResponse) EncodeJSON() gin.H {
	var result = make(gin.H)
	for k, v := range encryptedResponse.Params {
		result[k] = v
	}
	return result
}

func (encryptedResponse EncryptedResponse) decrypt(sessionKey SessionKey) (Response, error) {
	params, err := encryptedResponse.Params.Decrypt(sessionKey)
	return Response{Params: params}, err
}

func (response Response) encrypt(sessionKey SessionKey) EncryptedResponse {
	return EncryptedResponse{Params: response.Params.Encrypt(sessionKey)}
}
