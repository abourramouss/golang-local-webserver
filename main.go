package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func bla(w http.ResponseWriter, r *http.Request) string {
	http.FileServer(http.Dir("/Users/ayman/Desktop/")).ServeHTTP(w, r)

	return r.URL.Path
}

func handler(w http.ResponseWriter, r *http.Request, path string) {
	serr := r.ParseForm()

	if serr != nil {
		log.Fatal(serr)
	}
	name := r.PostFormValue("name")
	subpath := path + "/" + name
	err := os.Mkdir(subpath, 0755)

	fmt.Printf("%s\n", subpath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Carpeta creada correctamente en: %s", subpath)
}

func uploadFile(w http.ResponseWriter, r *http.Request, path string) {
	fmt.Fprint(w, "Uploading File")

	//1. parse input form, type multipart/form-data
	r.ParseMultipartForm(10 << 20)
	//2. retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error form-data")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//3. write temporary file on our server
	subpath := "/Users/ayman/Desktop" + path

	tempFile, err := ioutil.TempFile(subpath, "*-"+handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	//4. return wheter or not this has been successful

}
func setupRoutes() {
	var path string

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadFile(w, r, path)
		fmt.Fprintf(w, "Se subio el archivo correctamente en: %s", path)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("/Users/ayman/Desktop/")).ServeHTTP(w, r)
		path = r.URL.Path
		fmt.Printf("%s\n", path)

	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		path := "/Users/ayman/Desktop" + path
		handler(w, r, path)

	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Go File upload")
	setupRoutes()
}
