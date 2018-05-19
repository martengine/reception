package reception

// Service declaration.
type Service struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}
