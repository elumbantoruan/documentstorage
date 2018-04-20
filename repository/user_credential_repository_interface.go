package repository

import "github.com/elumbantoruan/documentstorage/model"

// UserCredentialRepository is an interface for user repository
type UserCredentialRepository interface {
	AddUser(creds model.Credentials) error
	GetUser(creds model.Credentials) error
}
