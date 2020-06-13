package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tempFolder string = "TEMP"
var formFileID string = "file"
var maxFileSize int64 = 10 << 20

//source: https://medium.com/@petehouston/upload-files-with-curl-93064dcccc76
func saveFile(r *http.Request) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(maxFileSize)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(formFileID)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//temp folder
	tempPath := filepath.Join(".", tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "upload-*-"+handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	return tempFile.Name(), nil
}

//ViewUploadFile save file
func ViewUploadFile(w http.ResponseWriter, r *http.Request) {
	log.Println("ViewUploadFile")
	fileName, err := saveFile(r)
	if err != nil {
		fmt.Fprint(w, "error")
		fmt.Println("File Upload: error", err)
		log.Println("File Upload: error ", err)
		return
	}

	// return that we have successfully uploaded our file!
	fmt.Printf("File Upload: uploaded file %s\n", fileName)
	log.Println("File Upload: uploaded file ", fileName)
	fmt.Fprintf(w, "ok")
}
