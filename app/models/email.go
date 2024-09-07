package models

type GithubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"` // Optional: sometimes GitHub includes visibility (e.g., "public")
}
