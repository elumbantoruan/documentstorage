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

var users = make(map[string]model.Credentials)

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

	if _, exists := existingUsers[strings.ToLower(creds.UserName)]; exists {
		err = fmt.Errorf("user %s exists already", creds.UserName)
		return err
	}

	var allUsers []model.Credentials

	for _, v := range existingUsers {
		allUsers = append(allUsers, v)
	}

	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.FileName)
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return err
	}

	bytes, _ := json.MarshalIndent(allUsers, "", "\t")
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
	// var mapper = make(map[string]string)
	// for _, e := range existingUsers {
	// 	mapper[strings.ToLower(e.UserName)] = e.Password
	// }
	credsInvalid := errors.New("username/login is not valid")
	if v, ok := existingUsers[strings.ToLower(creds.UserName)]; ok {
		if v.Password != creds.Password {
			return credsInvalid
		}
	} else {
		return credsInvalid
	}
	return nil
}

func (u *UserCredentialFileRepository) getAllUsers() (map[string]model.Credentials, error) {
	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.FileName)
	var creds []model.Credentials

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if len(bytes) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return nil, err
	}
	// lazy loading
	if len(users) == 0 {
		for _, eu := range creds {
			users[strings.ToLower(eu.UserName)] = eu
		}
	}

	return users, nil
}
