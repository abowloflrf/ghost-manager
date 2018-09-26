package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

//File 文件信息
type File struct {
	Name string
	Size float64
	Ext  string
}

var coverPath = ""
var uploadPath = ""

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Where is Ghost content path?")
	}
	coverPath = os.Args[1] + "/images/cover"
	uploadPath = os.Args[1] + "/upload"
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("./templates/index.html")

		files, _ := ioutil.ReadDir(uploadPath)
		fileInfos := []File{}
		for _, file := range files {
			if !file.IsDir() {
				// fmt.Println(file.Name())
				fileInfos = append(fileInfos, File{
					file.Name(),
					float64(file.Size()),
					path.Ext(file.Name()),
				})
			}
		}
		t.Execute(w, fileInfos)
	})

	r.HandleFunc("/upload/api/cover", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed"))
			return
		}
		w.Write([]byte("POST: change cover api"))
	})

	r.HandleFunc("/upload/api/attachment", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
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

		} else if req.Method == http.MethodDelete {
			w.Write([]byte("DELETE: attachment api"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed"))
			return
		}
	})

	r.PathPrefix("/upload/assets/").Handler(http.StripPrefix("/upload/assets/", http.FileServer(http.Dir("./assets"))))
	r.PathPrefix("/upload/data/").Handler(http.StripPrefix("/upload/data/", http.FileServer(http.Dir(uploadPath))))

	fmt.Println("Listen: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
