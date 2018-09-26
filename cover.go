package main

import (
	"io"
	"net/http"
	"os"
)

//ChangeCover change cover
func ChangeCover(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}

	//get data from request
	slug := req.FormValue("coverSlug")
	cover, _, err := req.FormFile("coverFile")
	if err != nil {
		panic(err)
	}
	defer cover.Close()

	//save file to our cover path
	f, err := os.OpenFile(coverPath+"/"+slug+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, cover)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"OK"}`))
}
