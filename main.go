package utapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Config for UploadThing
// Host is always https://api.uploadthing.com
// ApiKey is taken from the UPLOADTHING_SECRET environment variable
// Version - SDK version (e.g., 6.10.0)
type uploadthingConfig struct {
	Host    string
	ApiKey  string
	Version string
}

// Client structure
// httpClient can be overridden for tests
// config stores the key and version
// fePackage/beAdapter - optional headers for analytics (can be left empty)
type UtApi struct {
	config     *uploadthingConfig
	httpClient *http.Client
	fePackage  string
	beAdapter  string
}

// Loads environment variables from .env
func setEnvironmentVariablesFromFile() error {
	return godotenv.Load(".env")
}

func handleSetEnvironmentVariables() error {
	err := setEnvironmentVariablesFromFile()
	if err != nil {
		fmt.Printf("Failed to load environment variables from .env file\n")
		return err
	}
	return nil
}

func validateEnvironmentVariables(keys []string) error {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			return fmt.Errorf("%s is not set", key)
		}
	}
	return nil
}

func getUploadthingConfig() (*uploadthingConfig, error) {
	err := handleSetEnvironmentVariables()
	if err != nil {
		return nil, err
	}
	err = validateEnvironmentVariables([]string{"UPLOADTHING_SECRET"})
	if err != nil {
		return nil, err
	}
	return &uploadthingConfig{
		Host:    "https://api.uploadthing.com",
		ApiKey:  os.Getenv("UPLOADTHING_SECRET"),
		Version: "7.6.0",
	}, nil
}

// Client constructor
func NewUtApi() (*UtApi, error) {
	config, err := getUploadthingConfig()
	if err != nil {
		return nil, err
	}
	return &UtApi{
		config:     config,
		httpClient: &http.Client{},
		fePackage:  "",
		beAdapter:  "github.com/IXackerr/utapi-go",
	}, nil
}

// Set headers for all requests
func (ut *UtApi) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-uploadthing-api-key", ut.config.ApiKey)
	req.Header.Set("x-uploadthing-version", ut.config.Version)
	if ut.fePackage != "" {
		req.Header.Set("x-uploadthing-fe-package", ut.fePackage)
	}
	if ut.beAdapter != "" {
		req.Header.Set("x-uploadthing-be-adapter", ut.beAdapter)
	}
}

// Universal POST request
func (ut *UtApi) post(path string, body *bytes.Buffer) (*http.Response, error) {
	url := ut.config.Host + path
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	ut.setHeaders(req)
	resp, err := ut.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("UploadThing: error %d: %s", resp.StatusCode, string(respBody))
	}
	return resp, nil
}

// =====================
// ==== API Methods ====
// =====================

// 1. Delete files
// POST /v6/deleteFiles
// {"fileKeys": ["key1", ...]}
type DeleteFilesRequest struct {
	FileKeys []string `json:"fileKeys"`
}
type DeleteFilesResponse struct {
	Success      bool `json:"success"`
	DeletedCount int  `json:"deletedCount"`
}

func (ut *UtApi) DeleteFiles(fileKeys []string) (*DeleteFilesResponse, error) {
	payload := DeleteFilesRequest{FileKeys: fileKeys}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := ut.post("/v6/deleteFiles", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result DeleteFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 3. List files
// POST /v6/listFiles
// {"limit": 100, "offset": 0}
type ListFilesRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
type ListFilesFile struct {
	Id         string  `json:"id"`
	CustomId   *string `json:"customId"`
	Key        string  `json:"key"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	Size       int64   `json:"size"`
	UploadedAt int64   `json:"uploadedAt"`
}
type ListFilesResponse struct {
	HasMore bool            `json:"hasMore"`
	Files   []ListFilesFile `json:"files"`
}

func (ut *UtApi) ListFiles(limit, offset int) (*ListFilesResponse, error) {
	payload := ListFilesRequest{Limit: limit, Offset: offset}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := ut.post("/v6/listFiles", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ListFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 4. Rename files
// POST /v6/renameFiles
// {"updates": [{"fileKey": "...", "newName": "..."}]}
type RenameFileUpdate struct {
	FileKey string `json:"fileKey"`
	NewName string `json:"newName"`
}
type RenameFilesRequest struct {
	Updates []RenameFileUpdate `json:"updates"`
}
type RenameFilesResponse struct {
	Success      bool `json:"success"`
	RenamedCount int  `json:"renamedCount"`
}

func (ut *UtApi) RenameFiles(updates []RenameFileUpdate) (*RenameFilesResponse, error) {
	payload := RenameFilesRequest{Updates: updates}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := ut.post("/v6/renameFiles", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result RenameFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 5. Get usage info
// POST /v6/getUsageInfo
type UsageInfoResponse struct {
	TotalBytes    int64 `json:"totalBytes"`
	AppTotalBytes int64 `json:"appTotalBytes"`
	FilesUploaded int   `json:"filesUploaded"`
	LimitBytes    int64 `json:"limitBytes"`
}

func (ut *UtApi) GetUsageInfo() (*UsageInfoResponse, error) {
	resp, err := ut.post("/v6/getUsageInfo", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result UsageInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 6. Get presigned URL for a private file
// POST /v6/requestFileAccess
// {"fileKey": "...", "expiresIn": 3600}
type RequestFileAccessRequest struct {
	FileKey   string `json:"fileKey"`
	ExpiresIn int    `json:"expiresIn,omitempty"`
}
type RequestFileAccessResponse struct {
	UfsUrl string `json:"ufsUrl"`
	Url    string `json:"url"` // deprecated
}

func (ut *UtApi) GetPresignedUrl(fileKey string, expiresIn int) (string, error) {
	payload := RequestFileAccessRequest{FileKey: fileKey, ExpiresIn: expiresIn}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	resp, err := ut.post("/v6/requestFileAccess", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result RequestFileAccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.UfsUrl, nil
}

// 7. Get app info
// POST /v7/getAppInfo
type GetAppInfoResponse struct {
	AppId            string `json:"appId"`
	DefaultACL       string `json:"defaultACL"`
	AllowACLOverride bool   `json:"allowACLOverride"`
}

func (ut *UtApi) GetAppInfo() (*GetAppInfoResponse, error) {
	resp, err := ut.post("/v7/getAppInfo", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result GetAppInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 8. Get presigned URL for file upload (without file router)
// POST /v6/uploadFiles
// {"files": [{"name": "file.txt", "size": 123, "type": "text/plain"}], "acl": "public-read"}
type UploadFileInfo struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomId *string `json:"customId,omitempty"`
}
type UploadFilesRequest struct {
	Files              []UploadFileInfo `json:"files"`
	ACL                string           `json:"acl"` // "public-read" or "private"
	Metadata           interface{}      `json:"metadata,omitempty"`
	ContentDisposition string           `json:"contentDisposition,omitempty"`
}
type PresignedPostURLs struct {
	Key                string            `json:"key"`
	FileName           string            `json:"fileName"`
	FileType           string            `json:"fileType"`
	FileUrl            string            `json:"fileUrl"`
	ContentDisposition string            `json:"contentDisposition"`
	PollingJwt         string            `json:"pollingJwt"`
	PollingUrl         string            `json:"pollingUrl"`
	CustomId           *string           `json:"customId"`
	Url                string            `json:"url"`
	Fields             map[string]string `json:"fields"`
}
type UploadFilesResponse struct {
	Data []PresignedPostURLs `json:"data"`
}

// Get presigned URL for file upload
func (ut *UtApi) GetPresignedUploadUrl(files []UploadFileInfo, acl string) (*UploadFilesResponse, error) {
	payload := UploadFilesRequest{
		Files: files,
		ACL:   acl,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := ut.post("/v6/uploadFiles", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result UploadFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// createMultipartForm creates multipart/form-data body for S3 compatible POST
func createMultipartForm(content io.Reader, size int64, fileName string, fields map[string]string) (body *bytes.Buffer, contentType string, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, "", err
	}
	if content != nil && size > 0 {
		if _, err = io.CopyN(fw, content, size); err != nil {
			return nil, "", err
		}
	}
	w.Close()
	return &b, w.FormDataContentType(), nil
}

// Upload file to presigned URL (S3 compatible POST)
// filePath - path to local file
// presigned - PresignedPostURLs struct from GetPresignedUploadUrl response
func UploadFileToPresignedUrl(filePath string, presigned PresignedPostURLs) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	body, contentType, err := createMultipartForm(file, fileInfo.Size(), presigned.FileName, presigned.Fields)
	if err != nil {
		return err
	}

	// Send POST to presigned.Url
	req, err := http.NewRequest("POST", presigned.Url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("File upload error: %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// Upload content to presigned URL (S3 compatible POST)
// content - io.Reader with file data
// size - content size in bytes
// presigned - PresignedPostURLs struct from GetPresignedUploadUrl response
func UploadContentToPresignedUrl(content io.Reader, size int64, presigned PresignedPostURLs) error {
	body, contentType, err := createMultipartForm(content, size, presigned.FileName, presigned.Fields)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", presigned.Url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("File upload error: %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}