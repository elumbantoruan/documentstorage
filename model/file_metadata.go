package model

// FileStorage represent file storage information for specific user
type FileStorage struct {
	FileContent FileContent `json:"fileContent"`
	UserName    string      `json:"userName"`
	FileName    string      `json:"fileName"`
}

// FileStorageResponse represents a response for filestorage request
type FileStorageResponse struct {
	Location string `json:"location"`
}

// FileContent represents a response for request of a specific file
type FileContent struct {
	ContentLength int    `json:"contentLength"`
	ContentType   string `json:"contentType"`
	Payload       []byte `json:"payload"`
}
