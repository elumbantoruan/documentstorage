package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elumbantoruan/documentstorage/model"
	"github.com/elumbantoruan/documentstorage/repository"
	"github.com/gorilla/mux"
)

// Files handles user files storage
type Files struct {
	UserStorage repository.UserStorageRepository
}

// NewFiles returns an instance of Files Handler
func NewFiles(ufr repository.UserStorageRepository) *Files {
	return &Files{
		UserStorage: ufr,
	}
}

// HandleUploadFile handles upload file
func (f *Files) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userName, err := f.validateAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	fileName := ""

	if _, found := vars["fileName"]; found {
		fileName = vars["fileName"]
	} else {
		// no filename
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fc := model.FileContent{}
	err = json.NewDecoder(r.Body).Decode(&fc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}
	fr := model.FileStorage{
		FileContent: fc,
		UserName:    userName,
		FileName:    fileName,
	}
	resp, err := f.UserStorage.UploadFile(fr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// HandleGetFile handles get file
func (f *Files) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userName, err := f.validateAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println(userName)
}

// HandleGetFiles handles get all files for specific user
func (f *Files) HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userName, err := f.validateAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println(userName)
}

// HandleDeleteFile handles delete file
func (f *Files) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userName, err := f.validateAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println(userName)

}

func (f *Files) validateAuth(r *http.Request) (string, error) {
	auth := r.Header.Get("X-Session")
	if len(auth) == 0 {
		err := errors.New("empty x-session token")
		return "", err
	}

	token, _ := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if token == nil {
		err := errors.New("invalid token")
		return "", err
	}

	var userName string
	err := errors.New("unable to find username in the claim")

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		userName = claims["username"].(string)
		if len(userName) == 0 {
			return "", err
		}
	} else {
		return "", err
	}

	return userName, nil
}
