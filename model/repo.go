package model

type RepoInfo struct {
	ID        int64  `json:"id"`
	OwnerName string `json:"owner_name"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}
