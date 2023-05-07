package handlers

import (
	"log"
	"net/http"
	"html/template"
)

type Form struct {
	l *log.Logger
}

func NewForm(l *log.Logger) *Form {
	return &Form{l}
}

func (f *Form) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	f.l.Println("GET request")
	myTemplate := template.Must(template.ParseFiles("formData.html"))
	myTemplate.Execute(rw, nil)
}