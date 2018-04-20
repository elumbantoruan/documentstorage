package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/elumbantoruan/documentstorage/model"
	"github.com/elumbantoruan/documentstorage/repository"
)

// Register handles user registration
type Register struct {
	UserRepository repository.UserCredentialRepository
}

// NewRegistration returns an instance of Register
func NewRegistration(urp repository.UserCredentialRepository) *Register {
	return &Register{
		UserRepository: urp,
	}
}

// HandleRegister handles the registration
func (reg *Register) HandleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	creds := model.Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = reg.validateUserNameLength(creds.UserName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = reg.validateUserNameAlphaNumeric(creds.UserName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = reg.validatePasswordLength(creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = reg.validatePasswordCharacters(creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	err = reg.UserRepository.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.NewException(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (reg *Register) validateUserNameLength(s string) error {
	if len(s) < 3 || len(s) > 20 {
		return errors.New("user name must be at least 3 characters and no more than 20")
	}
	return nil
}

func (reg *Register) validateUserNameAlphaNumeric(s string) error {
	re := regexp.MustCompile("^[0-9A-Za-z]*$")
	valid := re.MatchString(s)
	if !valid {
		return errors.New("user name may only contain alphanumeric characters")
	}
	return nil
}

func (reg *Register) validatePasswordLength(s string) error {
	if len(s) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (reg *Register) validatePasswordCharacters(s string) error {
	re := regexp.MustCompile("^[A-Za-z]*$")
	valid := re.MatchString(s)
	if !valid {
		return errors.New("passwords must contain only characters")
	}
	return nil
}
