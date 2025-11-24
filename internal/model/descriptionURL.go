package model

type DescriptionURL struct {
	ID          int    `json:"uuid,omitempty"`
	Short       string `json:"short_url,omitempty"`
	Original    string `json:"original_url,omitempty"`
	Correlation string `json:"correlation_id,omitempty"`
}
