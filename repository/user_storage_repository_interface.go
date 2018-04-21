package repository

import "github.com/elumbantoruan/documentstorage/model"

// UserStorageRepository represents an interface for a User storage repository
type UserStorageRepository interface {
	UploadFile(request model.FileStorage) (*model.FileStorageResponse, error)
	GetFile(userName, fileName string) (*model.FileContent, error)
	GetFiles(userName string) ([]string, error)
	DeleteFile(userName, fileName string) (bool, error)
}
