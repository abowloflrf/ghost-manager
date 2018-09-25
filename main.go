package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

//File 文件信息
type File struct {
	Name string
	Size float64
	Ext  string
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, res *http.Request) {
		t, _ := template.ParseFiles("./templates/index.html")

		files, _ := ioutil.ReadDir("/Users/ruofeng/Downloads")
		fileInfos := []File{}
		for _, file := range files {
			if !file.IsDir() {
				fileInfos = append(fileInfos, File{
					file.Name(),
					float64(file.Size()) / 1048576,
					path.Ext(file.Name()),
				})
			}
		}
		t.Execute(w, fileInfos)
	})

	fmt.Println("Listen: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}