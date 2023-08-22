package go_printify

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	uploadsPath = "uploads.json"
	uploadPath  = "uploads/images.json"
)

type UploadsResponse struct {
	CurrentPage  int      `json:"current_page"`
	Data         []Upload `json:"data"`
	FirstPageUrl string   `json:"first_page_url"`
	LastPageUrl  string   `json:"last_page_url"`
	NextPageUrl  string   `json:"next_page_url"`
	From         int      `json:"from"`
	LastPage     int      `json:"last_page"`
	Path         string   `json:"path"`
	PerPage      int      `json:"per_page"`
	PrevPageUrl  string   `json:"prev_page_url"`
	To           int      `json:"to"`
	Total        int      `json:"total"`
}

type AddFileData struct {
	FileName string `json:"file_name"`
	Contents string `json:"contents"`
}

type AddFileUrl struct {
	FileName string `json:"file_name"`
	Contents string `json:"contents"`
}

type Upload struct {
	UploadId   string `json:"id"`
	FileName   string `json:"file_name"`
	Height     int    `json:"height"`
	Width      int    `json:"width"`
	Size       int    `json:"size"`
	MimeType   string `json:"mime_type"`
	PreviewUrl string `json:"preview_url"`
	UploadTime string `json:"upload_time"`
}

func (c *Client) GetAllUploads() ([]Upload, error) {

	var allUploads []Upload
	page := 1
	for {
		uploadResults, err := c.GetUploads(&page)
		if err != nil {
			fmt.Println("Received error from getUploads")
			return nil, err
		}

		allUploads = append(allUploads, uploadResults.Data...)

		if uploadResults.NextPageUrl == "" {
			break
		}
		page++
	}

	return allUploads, nil
}

func (c *Client) GetUploads(page *int) (*UploadsResponse, error) {

	query := fmt.Sprintf("limit=100&page=%d", *page)
	req, err := c.newRequest(http.MethodGet, uploadsPath, query, nil)
	if err != nil {
		return nil, err
	}
	uploads := &UploadsResponse{}
	_, err = c.do(req, uploads)
	return uploads, err
}

func (c *Client) AddUpload(path string) (*Upload, error) {

	// Get filename
	filename := filepath.Base(path)

	// Validate 5Mb or less
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	fileSize := info.Size()
	if fileSize > 5*1024*1024 {
		return nil, fmt.Errorf("file should not be more than 5MB")
	}

	// Read file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// base64 encode file
	encoded := base64.StdEncoding.EncodeToString(contents)

	body := AddFileData{
		FileName: filename,
		Contents: encoded,
	}

	// Send upload
	req, err := c.newRequest(http.MethodPost, uploadPath, "", body)
	if err != nil {
		return nil, err
	}
	upload := &Upload{}
	_, err = c.do(req, upload)
	return upload, err
}

func (c *Client) AddUploads(uploadItems []string) (*Upload, error) {
	for _, element := range uploadItems  {
		c.AddUpload(element)
	}
	return nil, nil
}