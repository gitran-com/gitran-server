package model

type RepoInfo struct {
	ID        uint64 `json:"id"`
	OwnerName string `json:"owner_name"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}
