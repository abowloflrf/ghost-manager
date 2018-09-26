package main

import (
	"encoding/json"
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

	r.HandleFunc("/upload/api/upload-attachment", func(w http.ResponseWriter, req *http.Request) {
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
		// w.Write([]byte(`{"status":"OK"}`))
	})

	r.HandleFunc("/upload/api/delete-attachment", func(w http.ResponseWriter, req *http.Request) {
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
			w.Write([]byte(`{"status":"ERROR","msg":"File Not Found"}`))
			return
		}
		//delete file
		err = os.Remove(uploadPath + "/" + fileName)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"status":"OK"}`)
		// w.Write([]byte(`{"status":"OK"}`))

	})

	//serve webiste's static files: css,js,img...
	r.PathPrefix("/upload/assets/").Handler(http.StripPrefix("/upload/assets/", http.FileServer(http.Dir("./assets"))))
	//serve user upload attachments in Ghost content directory
	r.PathPrefix("/upload/data/").Handler(http.StripPrefix("/upload/data/", http.FileServer(http.Dir(uploadPath))))

	fmt.Println("Listen: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
