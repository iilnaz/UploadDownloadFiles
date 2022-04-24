package controller

import (
	"FileUploadDownload/utils"
	"fmt"
	"net/http"
	"path"
)

const (
	DownloadPath = "download"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, fmt.Sprintf("Uploaded File: %+v\n", header.Filename))
	fmt.Fprintf(w, fmt.Sprintf("File Size: %+v\n", header.Size))
	fmt.Fprintf(w, fmt.Sprintf("MIME Header: %+v\n", header.Header))

	err = utils.ChunkFile(file, header)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uploadFile(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown Method. Try to do POST\n")
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := path.Base(r.URL.Path)
	if fileName == DownloadPath {
		fmt.Fprintf(w, "Write the file name\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := utils.GlueFiles(w, fileName); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		downloadFile(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown Method. Try to do GET\n")
	}
}
