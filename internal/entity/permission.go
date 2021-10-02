package entity

// Permission represents a permission.
type Permission struct {
	Rules       []string `json:"rules"`
	SubjectName string   `json:"subject_name"`
}
