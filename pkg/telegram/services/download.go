package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetFile devuelve el file_path interno de Telegram para un file_id: hace
// falta para poder descargar el archivo real.
func (s *Services) GetFile(ctx context.Context, fileID string) (string, error) {
	params := url.Values{}
	params.Set("file_id", fileID)

	var result struct {
		FilePath string `json:"file_path"`
	}
	if err := s.call(ctx, "getFile", params, &result); err != nil {
		return "", err
	}
	return result.FilePath, nil
}

// DownloadPhoto descarga el contenido real de una foto a partir de su
// file_id.
func (s *Services) DownloadPhoto(ctx context.Context, fileID string) ([]byte, error) {
	filePath, err := s.GetFile(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return s.DownloadFile(ctx, filePath)
}

// DownloadFile descarga el contenido real de un archivo de Telegram, dado
// el file_path devuelto por GetFile. Los file_id/file_path de Telegram no
// son URLs permanentes: hay que descargar y volver a subir a donde se
// quiera guardar el archivo (ver .claude/rules/notion-rules.md).
func (s *Services) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/file/bot%s/%s", s.baseURL, s.token, filePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("telegram: build download request: %w", withoutURL(err))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("telegram: download failed: %w", withoutURL(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("telegram: unexpected status %d downloading file", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("telegram: read downloaded file: %w", err)
	}
	return data, nil
}
