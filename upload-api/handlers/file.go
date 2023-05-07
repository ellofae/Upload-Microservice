package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ellofae/Upload-Microservice/upload-api/data"
	"github.com/gorilla/mux"
)

type FileHandler struct {
	l     *log.Logger
	store data.Storage
}

func NewFileHandler(l *log.Logger, s data.Storage) *FileHandler {
	return &FileHandler{l, s}
}

func (f *FileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	f.l.Println("POST request")

	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	if id == "" || fn == "" {
		f.l.Println("Bad request: mux variables are not in the correct format")
		http.Error(rw, "Didn't manage to get id or filename from your request", http.StatusBadRequest)
		return
	}

	f.saveFile(id, fn, rw, r)
}

func (f *FileHandler) saveFile(id string, path string, rw http.ResponseWriter, r *http.Request) {
	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		http.Error(rw, "Didn't manage to save the file", http.StatusInternalServerError)
		return
	}
}
