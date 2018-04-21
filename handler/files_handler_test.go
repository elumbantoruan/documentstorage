package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/elumbantoruan/documentstorage/model"
	"github.com/elumbantoruan/documentstorage/repository"
	"github.com/stretchr/testify/assert"
)

func TestFiles_HandleUploadFile_MissingXToken(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("PUT", "files/test1.txt", strings.NewReader(string(bytes)))
	responseRecorder := httptest.NewRecorder()

	// Register files resource
	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleUploadFile(responseRecorder, request)

	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestFiles_HandleUploadFile_Success(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("PUT", "files/test1.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test1.txt",
	}
	request = mux.SetURLVars(request, mapper)
	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	// Register files resource
	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleUploadFile(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
}

func TestFiles_HandleGetFile_UserExistsFileNotFound(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("GET", "files/test111.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test111.txt",
	}
	request = mux.SetURLVars(request, mapper)

	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	// Register files resource
	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleGetFile(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestFiles_HandleGetFile_UserExistsFileFound(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("GET", "files/test4.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test4.txt",
	}
	request = mux.SetURLVars(request, mapper)

	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	// Register files resource
	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleGetFile(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	fc := model.FileContent{}
	err := json.NewDecoder(responseRecorder.Body).Decode(&fc)
	assert.NoError(t, err)
	assert.NotNil(t, fc)
}

func TestFiles_HandleGetFiles_NotFound(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("GET", "files", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test111.txt",
	}
	request = mux.SetURLVars(request, mapper)

	login := Login{}
	token, _ := login.CreateToken("user5")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleGetFiles(responseRecorder, request)
	// user5 will not be found
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestFiles_HandleGetFiles_Found(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("GET", "files", strings.NewReader(string(bytes)))

	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleGetFiles(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestFiles_HandleDeleteFiles_Found(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("DELETE", "files/test4.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test4.txt",
	}
	request = mux.SetURLVars(request, mapper)

	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleDeleteFile(responseRecorder, request)
	// user4 exists and test4.txt is deleted
	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
}

func TestFiles_HandleDeleteFiles_NotFound(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("DELETE", "files/test5.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test5.txt",
	}
	request = mux.SetURLVars(request, mapper)

	login := Login{}
	token, _ := login.CreateToken("user4")
	request.Header.Add("x-session", token)

	responseRecorder := httptest.NewRecorder()

	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleDeleteFile(responseRecorder, request)
	// user4 exists but test5.txt is not found
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestFiles_HandleDeleteFiles_NotLogin(t *testing.T) {

	payload := model.FileContent{
		ContentType:   "text",
		ContentLength: 20,
	}
	bytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("DELETE", "files/test4.txt", strings.NewReader(string(bytes)))
	mapper := map[string]string{
		"fileName": "test4.txt",
	}
	request = mux.SetURLVars(request, mapper)

	responseRecorder := httptest.NewRecorder()

	fs := repository.NewUserStorageMockRepository()
	stor := NewFiles(fs)
	stor.HandleDeleteFile(responseRecorder, request)
	// no token
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}
