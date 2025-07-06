package models

type ProcessingResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	ZipPath     string   `json:"zip_path,omitempty"`
	DownloadURL string   `json:"download_url,omitempty"`
	FrameCount  int      `json:"frame_count,omitempty"`
	Images      []string `json:"images,omitempty"`
}
