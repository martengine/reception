package service

// Schema is a list of data kept by a service.
type Schema map[string]Data

// Data represents information that might be contained in Service Schema.
type Data struct {
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

// Service declaration.
type Service struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}
