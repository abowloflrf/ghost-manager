package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

//UploadAttachment attachment upload handler
func UploadAttachment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	//get file from request
	file, header, err := req.FormFile("attachment")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//TODO: check if there exsit file with same name
	//save file to our path
	f, err := os.OpenFile(uploadPath+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, file)

	//return json
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"status":"OK"}`)
}

//DeleteAttachment attachment deletion handler
func DeleteAttachment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}

	reqData := struct {
		Filename string
	}{}
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		panic(err)
	}
	fileName := reqData.Filename
	//check file exist
	_, err = os.Stat(uploadPath + "/" + fileName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"status":"ERROR","msg":"File Not Found"}`)
		return
	}
	//delete file
	err = os.Remove(uploadPath + "/" + fileName)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"status":"OK"}`)
}
