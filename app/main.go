package main

import (
	"FileUploadDownload/controller"
	"net/http"
	"os"
)

func init() {
	err := os.MkdirAll("storage", 0777)
	if err != nil {
		panic(err)
	}
}
func main() {
	http.HandleFunc("/upload", controller.UploadHandler)
	http.HandleFunc("/download/", controller.DownloadHandler)
	http.ListenAndServe(":8080", nil)
}
