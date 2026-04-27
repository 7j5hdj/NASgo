package main

import (
	dat "nasgo/progect/data"
	"nasgo/progect/handle"
	"nasgo/progect/scr"

	"log"
	"net/http"
	"os"
)

func main() {
	if _, err := os.Stat(dat.STORAGE_DIR); os.IsNotExist(err) {
		os.Mkdir(dat.STORAGE_DIR, 0755)
	}

	storage := scr.NewStorage()

	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		handle.HandleUpload(w, r, storage)
	})
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		handle.HandleDownload(w, r, storage)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		handle.HandleList(w, r, storage)
	})

	http.HandleFunc("/commands", handle.HandleCommands)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
