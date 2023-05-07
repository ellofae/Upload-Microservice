package handlers

import (
	"log"
	"net/http"
	"strconv"
	"io"
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

func (f *FileHandler) UploadREST(rw http.ResponseWriter, r *http.Request) {
	f.l.Println("POST request")

	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	if id == "" || fn == "" {
		f.l.Println("Bad request: mux variables are not in the correct format")
		http.Error(rw, "Didn't manage to get id or filename from your request", http.StatusBadRequest)
		return
	}

	f.saveFile(id, fn, rw, r.Body)
}

func (f *FileHandler) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128*1024)
	if err != nil {
		f.l.Println("Bad request", err)
		http.Error(rw, "Expected multipart form-data", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.l.Println("Process form for id", id)
	if idErr != nil {
		f.l.Println("Bad request", err)
		http.Error(rw, "Expected integer id", http.StatusBadRequest)
		return
	}

	fl, mh, err := r.FormFile("file")
	if err != nil {
		f.l.Println("Bad request", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(r.FormValue("id"), mh.Filename, rw, fl)
}

func (f *FileHandler) saveFile(id string, path string, rw http.ResponseWriter, r io.Reader) {
	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		http.Error(rw, "Didn't manage to save the file", http.StatusInternalServerError)
		return
	}
}
