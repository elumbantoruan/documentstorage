package handler

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elumbantoruan/documentstorage/model"
	"github.com/elumbantoruan/documentstorage/repository"
)

// Login handles user login
type Login struct {
	UserRepository repository.UserCredentialRepository
}

// NewLogin returns an instance of Login Handler
func NewLogin(urp repository.UserCredentialRepository) *Login {
	return &Login{
		UserRepository: urp,
	}
}

// HandleLogin handles user login
func (login *Login) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	creds := model.Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = login.UserRepository.GetUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	token, _ := login.createToken(creds.UserName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.NewUserToken(token))
}

func (login *Login) createToken(userName string) (string, error) {
	signKey := []byte("mysupersecretkey")
	claims := &jwt.MapClaims{
		"username": userName,
		"exp":      time.Now().Add(1 * time.Hour).UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	return ss, err
}
