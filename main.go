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

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Where is cover path?")
	}
	coverPath = os.Args[1]
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", func(w http.ResponseWriter, res *http.Request) {
		t, _ := template.ParseFiles("./templates/index.html")

		files, _ := ioutil.ReadDir(coverPath)
		fileInfos := []File{}
		for _, file := range files {
			if !file.IsDir() {
				// fmt.Println(file.Name())
				fileInfos = append(fileInfos, File{
					file.Name(),
					float64(file.Size()) / 1048576,
					path.Ext(file.Name()),
				})
			}
		}
		t.Execute(w, fileInfos)
	})

	r.PathPrefix("/upload/assets/").Handler(http.StripPrefix("/upload/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Println("Listen: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
