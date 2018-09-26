package main

import (
	"fmt"
	"html/template"
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

	//render index.html page
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

	//HTTP api
	r.HandleFunc("/upload/api/cover", ChangeCover)
	r.HandleFunc("/upload/api/upload-attachment", UploadAttachment)
	r.HandleFunc("/upload/api/delete-attachment", DeleteAttachment)

	//serve webiste's static files: css,js,img...
	r.PathPrefix("/upload/assets/").Handler(http.StripPrefix("/upload/assets/", http.FileServer(http.Dir("./assets"))))
	//serve user upload attachments in Ghost content directory
	r.PathPrefix("/upload/data/").Handler(http.StripPrefix("/upload/data/", http.FileServer(http.Dir(uploadPath))))

	fmt.Println("Listen: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
