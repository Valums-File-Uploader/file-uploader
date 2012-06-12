package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(fileName string, r io.Reader) error {
	saveToFilePath := filepath.Join("d:/tmp", fileName)
	osFile, err := os.Create(saveToFilePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, r)
	if err != nil {
		return err
	}
	return nil
}

func GetFileName(fh *multipart.FileHeader) (string, error) {
	itemHead := fh.Header["Content-Disposition"][0]
	lookfor := `filename="`
	fileIndex := strings.Index(itemHead, lookfor)
	if fileIndex < 0 {
		return "", errors.New(fmt.Sprintf("Can't find the file name: %s", itemHead))
	}
	//remove the last char "
	filePath := itemHead[fileIndex+len(lookfor) : len(itemHead)-1]
	_, fileName := filepath.Split(filePath)
	return fileName, nil
}

func UploadHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Recevied upload request")
	
	var reader io.Reader
	var fileName string

	if formFile, fh, err := req.FormFile("qqfile"); err == nil {
		fileName, err = GetFileName(fh)
		if err != nil {
			log.Println("GetFileName, err=", err)
			w.Write([]byte(`{"error":"Can NOT find the file name"}`))
			return
		}
		reader = formFile

	}else if req.FormValue("qqfile") != "" {
		fileName = req.FormValue("qqfile")
		reader = req.Body

	}else{
		w.Write([]byte(`{"error":"Can NOT find the file name"}`))
		return
	}

	log.Println("FileName =", fileName)
	err := SaveFile(fileName, reader)
	if err != nil {
		w.Write([]byte(`{"error":"Can NOT save the file"}`))
		log.Println("SaveFile ERR =", err)
		return
	}
	//all is success
	log.Println(fileName, "is saved")
	w.Write([]byte(`{"success":true}`))
	return
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":8080", nil)
}
