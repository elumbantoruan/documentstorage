package repository

import (
	"strings"

	"github.com/elumbantoruan/documentstorage/model"
)

// MockUserStorageRepository implements UserStorageRepository interface in file system
type MockUserStorageRepository struct {
	Files []model.FileStorage
}

// NewUserStorageMockRepository returns an instance of MockUserStorageRepository
func NewUserStorageMockRepository() *MockUserStorageRepository {

	return &MockUserStorageRepository{}
}

// UploadFile upload user file into UserStorageRepository
func (u *MockUserStorageRepository) UploadFile(request model.FileStorage) (*model.FileStorageResponse, error) {
	fs, err := u.getFiles()
	if err != nil {
		return nil, err
	}
	fs = append(fs, request)

	fsr := model.FileStorageResponse{
		Location: request.FileName,
	}
	return &fsr, nil
}

// GetFile returns a user file
func (u *MockUserStorageRepository) GetFile(userName, fileName string) (*model.FileStorage, error) {
	allFiles, err := u.getFiles()
	if err != nil {
		return nil, err
	}
	for _, f := range allFiles {
		if strings.EqualFold(f.UserName, userName) &&
			strings.EqualFold(f.FileName, fileName) {
			return &f, nil
		}
	}
	return nil, nil
}

// DeleteFile delete user file from storage
func (u *MockUserStorageRepository) DeleteFile(userName, fileName string) (bool, error) {
	files, err := u.getFiles()
	if err != nil {
		return false, err
	}
	found := false
	for i, f := range files {
		if strings.EqualFold(f.UserName, userName) &&
			strings.EqualFold(f.FileName, fileName) {
			files = append(files[:i], files[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	return true, nil
}

// GetFiles returns all files belong to a specific user
func (u *MockUserStorageRepository) GetFiles(userName string) ([]string, error) {
	allFiles, err := u.getFiles()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range allFiles {
		if strings.EqualFold(f.UserName, userName) {
			files = append(files, f.FileName)
		}
	}
	return files, nil
}

// GetFiles returns all files belong to a specific user
func (u *MockUserStorageRepository) getFiles() ([]model.FileStorage, error) {
	// var files []model.FileStorage
	u.Files = append(u.Files, model.FileStorage{
		FileName: "test1.txt",
		UserName: "user1",
		FileContent: model.FileContent{
			ContentLength: 10,
			ContentType:   "text",
			Payload:       nil,
		},
	})
	u.Files = append(u.Files, model.FileStorage{
		FileName: "test2.txt",
		UserName: "user2",
		FileContent: model.FileContent{
			ContentLength: 10,
			ContentType:   "text",
			Payload:       nil,
		},
	})
	u.Files = append(u.Files, model.FileStorage{
		FileName: "test3.txt",
		UserName: "user3",
		FileContent: model.FileContent{
			ContentLength: 10,
			ContentType:   "text",
			Payload:       nil,
		},
	})
	u.Files = append(u.Files, model.FileStorage{
		FileName: "test4.txt",
		UserName: "user4",
		FileContent: model.FileContent{
			ContentLength: 10,
			ContentType:   "text",
			Payload:       nil,
		},
	})
	return u.Files, nil
}
