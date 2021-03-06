package main

import (
	"log"
	"net/http"

	"github.com/elumbantoruan/documentstorage/handler"
	"github.com/elumbantoruan/documentstorage/repository"
	"github.com/gorilla/mux"
)

func main() {
	m, err := registerHandlers()
	if err != nil {
		log.Println(err)
		return
	}
	http.Handle("/", m)

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Println(err)
	}
}

func registerHandlers() (*mux.Router, error) {
	m := mux.NewRouter()

	// Register register resource
	ucfile := "/credentials/creds.json"
	uc := repository.NewUserCredentialFileRepository(ucfile)
	reg := handler.NewRegistration(uc)
	m.HandleFunc("/register", reg.HandleRegister).Methods("POST")

	// Register login resource
	login := handler.NewLogin(uc)
	m.HandleFunc("/login", login.HandleLogin).Methods("POST")

	// Register files resource
	fsfile := "/storage/storage.json"
	fs := repository.NewUserStorageFileRepository(fsfile)
	stor := handler.NewFiles(fs)
	m.HandleFunc("/files/{fileName}", stor.HandleUploadFile).Methods("PUT")
	m.HandleFunc("/files/{fileName}", stor.HandleGetFile).Methods("GET")
	m.HandleFunc("/files", stor.HandleGetFiles).Methods("GET")
	m.HandleFunc("/files/{fileName}", stor.HandleDeleteFile).Methods("DELETE")

	return m, nil
}
