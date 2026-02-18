package model

// DescriptionURL - структура для описания URL
// generate:reset
type DescriptionURL struct {
	ID          int    `json:"uuid,omitempty" db:"uuid,omitempty"`
	Short       string `json:"short_url,omitempty" db:"short_url,omitempty"`
	Original    string `json:"original_url,omitempty" db:"original_url,omitempty"`
	Correlation string `json:"correlation_id,omitempty"`
	UserID      string `json:"user_id,omitempty" db:"user_id,omitempty"`
	IsDeleted   *bool  `json:"is_deleted,omitempty" db:"is_deleted,omitempty"`
}

func (d *DescriptionURL) Copy() *DescriptionURL {
	if d == nil {
		return nil
	}
	return &DescriptionURL{
		ID:          d.ID,
		Short:       d.Short,
		Original:    d.Original,
		Correlation: d.Correlation,
		UserID:      d.UserID,
	}
}
