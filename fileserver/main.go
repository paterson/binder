package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paterson/binder/utils/logger"
	"github.com/paterson/binder/utils/request"
	"io"
	"mime/multipart"
	"os"
)

func main() {
	router := gin.Default()
	router.POST("/read", read)
	router.POST("/write", write)
	router.Run(port())
}

func read(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	fmt.Println(fmt.Sprintf("%+v", req.Params))
	if err == nil {
		filepath := req.Params["filepath"]
		req.SendFile("./.files" + filepath) // e.g ./.files/test.png
	}
}

func write(ctx *gin.Context) {
	req, err := request.Authenticate(ctx)
	if err == nil { // Auth is valid
		file, filename, err := req.RetrieveUploadedFile()
		checkError(err)
		err = storeFile(file, filename)
		checkError(err)
		req.Respond(request.StatusOK, request.Params{"success": "true"})
	}
}

func storeFile(file multipart.File, filename string) error {
	os.Mkdir("./.files/", 0777)
	out, err := os.Create("./.files/" + filename)
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}

func port() string {
	return ":" + os.Args[1]
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
