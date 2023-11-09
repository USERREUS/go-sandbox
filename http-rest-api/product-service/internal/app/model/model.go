package model

type Product struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}
