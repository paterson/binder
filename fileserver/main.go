package fileserver

import (
    "io"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/paterson/binder/authservice"
    "mime/multipart"
)

func main() {
    router := gin.Default()
    router.POST("/write", write)
    router.Static("/", "./.files")
    router.Run(":3001")
}

func write(ctx *gin.Context) {
    request, err := authservice.Authenticate(ctx)
    if err == nil { // Auth is valid
        file, filename, err := request.RetrieveUploadedFile()
        checkError(err)
        err = storeFile(file, filename)
        checkError(err)
        request.Respond(http.StatusOK, Body{"success": "true"})
    }
}

func storeFile(file multipart.File, filename string) error {
    out, err := os.Create("./.files/"+filename)
    defer out.Close()
    _, err = io.Copy(out, file)
    return err
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}
