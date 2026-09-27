package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// Upload sube un archivo al File Upload API de Notion en dos pasos (crear +
// enviar contenido) y devuelve el ID del file upload creado, listo para
// referenciar en la propiedad "Receipt" de una página (ver
// mapTransactionToProperties). Implementa ports.ReceiptUploader.
func (s *Services) Upload(data []byte, filename, contentType string) (string, error) {
	uploadID, uploadURL, err := s.createFileUpload()
	if err != nil {
		return "", err
	}

	if err := s.sendFileUpload(uploadURL, data, filename, contentType); err != nil {
		return "", err
	}

	return uploadID, nil
}

func (s *Services) createFileUpload() (id, uploadURL string, err error) {
	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/file_uploads", bytes.NewReader([]byte("{}")))
	if err != nil {
		return "", "", fmt.Errorf("notion: build file upload create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("notion: file upload create failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("notion: read file upload create response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("notion: unexpected status %d creating file upload: %s", resp.StatusCode, body)
	}

	var created struct {
		ID        string `json:"id"`
		UploadURL string `json:"upload_url"`
	}
	if err := json.Unmarshal(body, &created); err != nil {
		return "", "", fmt.Errorf("notion: decode file upload create response: %w", err)
	}
	return created.ID, created.UploadURL, nil
}

func (s *Services) sendFileUpload(uploadURL string, data []byte, filename, contentType string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, filename))
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("notion: build multipart file field: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return fmt.Errorf("notion: write file bytes: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("notion: close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, &body)
	if err != nil {
		return fmt.Errorf("notion: build file upload send request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("notion: file upload send failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("notion: read file upload send response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notion: unexpected status %d sending file upload: %s", resp.StatusCode, respBody)
	}
	return nil
}
