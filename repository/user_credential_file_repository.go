package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elumbantoruan/documentstorage/model"
)

// UserCredentialFileRepository implements UserRepository interface in file system
type UserCredentialFileRepository struct {
	FileName    string
	currentPath string
}

// NewUserCredentialFileRepository returns an instance of UserFileRepository
func NewUserCredentialFileRepository(filename string) *UserCredentialFileRepository {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return &UserCredentialFileRepository{
		FileName:    filename,
		currentPath: dir,
	}
}

// AddUser adds user into UserRepository
func (u *UserCredentialFileRepository) AddUser(creds model.Credentials) error {

	existingUsers, err := u.getAllUsers()
	if err != nil {
		return err
	}

	var users = make(map[string]interface{})
	for _, eu := range existingUsers {
		users[strings.ToLower(eu.UserName)] = nil
	}
	if _, exists := users[strings.ToLower(creds.UserName)]; exists {
		err = fmt.Errorf("user %s exists already", creds.UserName)
		return err
	}

	existingUsers = append(existingUsers, creds)

	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.FileName)
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return err
	}

	bytes, _ := json.MarshalIndent(existingUsers, "", "\t")
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns a token if user exists
func (u *UserCredentialFileRepository) GetUser(creds model.Credentials) error {
	existingUsers, err := u.getAllUsers()
	if err != nil {
		return err
	}
	var mapper = make(map[string]string)
	for _, e := range existingUsers {
		mapper[strings.ToLower(e.UserName)] = e.Password
	}
	credsInvalid := errors.New("username/login is not valid")
	if v, ok := mapper[strings.ToLower(creds.UserName)]; ok {
		if v != creds.Password {
			return credsInvalid
		}
	} else {
		return credsInvalid
	}
	return nil
}

func (u *UserCredentialFileRepository) getAllUsers() ([]model.Credentials, error) {
	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.FileName)
	var creds []model.Credentials

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if len(bytes) == 0 {
		return creds, nil
	}
	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return nil, err
	}
	return creds, nil
}
