package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elumbantoruan/documentstorage/model"
)

// UserStorageFileRepository implements UserStorageRepository interface in file system
type UserStorageFileRepository struct {
	StorageFileName string
	currentPath     string
}

// NewUserStorageFileRepository returns an instance of UserFileRepository
func NewUserStorageFileRepository(filename string) *UserStorageFileRepository {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return &UserStorageFileRepository{
		StorageFileName: filename,
		currentPath:     dir,
	}
}

// UploadFile upload user file into UserStorageRepository
func (u *UserStorageFileRepository) UploadFile(request model.FileStorage) (*model.FileStorageResponse, error) {
	fs, err := u.getFiles()
	if err != nil {
		return nil, err
	}
	fs = append(fs, request)

	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.StorageFileName)
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	bytes, _ := json.MarshalIndent(fs, "", "\t")
	_, err = f.Write(bytes)
	if err != nil {
		return nil, err
	}
	fsr := model.FileStorageResponse{
		Location: request.FileName,
	}
	return &fsr, nil
}

// GetFile returns a user file
func (u *UserStorageFileRepository) GetFile(userName, fileName string) (*model.FileContent, error) {
	allFiles, err := u.getFiles()
	if err != nil {
		return nil, err
	}
	for _, f := range allFiles {
		if strings.EqualFold(f.UserName, userName) &&
			strings.EqualFold(f.FileName, fileName) {
			return &f.FileContent, nil
		}
	}
	return nil, nil
}

// DeleteFile delete user file from storage
func (u *UserStorageFileRepository) DeleteFile(userName, fileName string) (bool, error) {
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
	sfname := fmt.Sprintf("%s/%s", u.currentPath, u.StorageFileName)
	f, err := os.OpenFile(sfname, os.O_TRUNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return false, err
	}

	bytes, _ := json.MarshalIndent(files, "", "\t")
	_, err = f.Write(bytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetFiles returns all files belong to a specific user
func (u *UserStorageFileRepository) GetFiles(userName string) ([]string, error) {
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
func (u *UserStorageFileRepository) getFiles() ([]model.FileStorage, error) {
	fileName := fmt.Sprintf("%s/%s", u.currentPath, u.StorageFileName)
	var files []model.FileStorage

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if len(bytes) == 0 {
		return files, nil
	}
	err = json.Unmarshal(bytes, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
