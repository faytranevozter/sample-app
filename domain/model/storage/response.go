package storage_model

type UploadResponse struct {
	Key         string `json:"key"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}
