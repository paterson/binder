package request

import (
	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/utils/encryption"
)

type Params map[string]string
type EncryptedParams map[string]string

func NewEncryptedParams(ctx *gin.Context) EncryptedParams {
	var dict map[string][]string
	if ctx.Request.Method == "POST" {
		ctx.Request.ParseMultipartForm(32 << 20)
		dict = ctx.Request.PostForm
	} else {
		dict = ctx.Request.URL.Query()
	}
	var params = make(EncryptedParams)
	for k, v := range dict {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}

func (params EncryptedParams) Decrypt(sessionKey SessionKey) (Params, error) {
	var result = make(Params)
	for k, v := range params {
		if k == "ticket" {
			result[k] = v
		} else {
			val, err := encryption.Decrypt(string(sessionKey), v)
			result[k] = val
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func (params Params) Encrypt(sessionKey SessionKey) EncryptedParams {
	var result = make(EncryptedParams)
	for k, v := range params {
		if k == "ticket" {
			result[k] = v
		} else {
			result[k] = encryption.Encrypt(string(sessionKey), v)
		}

	}
	return result
}
