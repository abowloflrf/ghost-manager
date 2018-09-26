package main

import "net/http"

//ChangeCover change cover
func ChangeCover(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	w.Write([]byte("POST: change cover api"))
}
