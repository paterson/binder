package clientproxy

import (
    "github.com/paterson/binder/utils/request"
)

const (
    signupURL       = "http://localhost:3000/signup"
    loginURL        = "http://localhost:3000/login"
    readRequestURL  = "http://localhost:3003/request/read"
    writeRequestURL = "http://localhost:3003/request/write"
)

type ClientProxy struct {
    Ticket     string
    SessionKey string
    http       gorequest
}

type AuthenticateParams struct {
    username  string
    password  string
}

type FileParams struct {
    ticket     string
    sessionKey string
    filepath   string
}

func NewClientPolicy() *ClientProxy {
    return ClientProxy{
        http: gorequest.New(),
    }
}

func (clientProxy *ClientProxy) Signup(username string, password string)  {
    params := AuthenticateParams{username: username, password: password}
    resp, body, errs := clientProxy.http.Post(signupURL).Send(params).End()
}

func (clientProxy *ClientProxy) Login(username string, password string) {
    params := AuthenticateParams{username: username, password: password}
    resp, body, errs := clientProxy.http.Post(loginURL).Send(params).End()
}

func (clientProxy *ClientProxy) Read(filepath) {
    params := FileParams{
        ticket:     clientProxy.Ticket,
        sessionKey: clientProxy.SessionKey,
        filepath:   filepath
    }
    resp, body, errs := clientProxy.http.Post(readRequestURL).Send(params).End()
    var dat map[string]interface{}
    err := json.Unmarshal(body, &dat)
    if err == nil {
        return
    }
    resp, body, errs := clientProxy.http.Post(dat["host"]).Send(params).End()

}

func (clientProxy *ClientProxy) Write(filepath, *File) {
    params := FileParams{
        ticket:     clientProxy.Ticket,
        sessionKey: clientProxy.SessionKey,
        filepath:   filepath
    }
    resp, body, errs := clientProxy.http.Post(writeRequestURL).Send(params).End()
    var dat map[string]interface{}
    err := json.Unmarshal(body, &dat)
    if err == nil {
        return
    }
    // Add file
    resp, body, errs := clientProxy.http.Post(dat["host"]).Send(params).End()
}

func
