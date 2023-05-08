package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ellofae/Upload-Microservice/upload-api/data"
	"github.com/ellofae/Upload-Microservice/upload-api/handlers"
	"github.com/gorilla/mux"
	gohandlers "github.com/gorilla/handlers"
)

const basePath string = "./filestorage"

func main() {
	l := log.New(os.Stdout, "upload-api", log.LstdFlags)

	localStorage, err := data.NewLocal(l, basePath, 1024)
	if err != nil {
		l.Fatal(err)
	}

	fileHandler := handlers.NewFileHandler(l, localStorage)
	mw := handlers.GzipHandler{}

	sm := mux.NewRouter()

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/files/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fileHandler.UploadREST)
	postRouter.HandleFunc("/", fileHandler.UploadMultipart)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/files/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/files/", http.FileServer(http.Dir(basePath))))
	getRouter.Use(mw.GzipMiddleware)

	// CORS set up
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))


	srv := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recived terminate, gracefil shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)
}
