package service

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateProjRequest struct {
	Name     string   `json:"name"`
	URI      string   `json:"uri"`
	Desc     string   `json:"desc"`
	GitURL   string   `json:"git_url"`
	SrcLangs []string `json:"src_langs"`
	TrnLangs []string `json:"trn_langs"`
	Type     int      `json:"type"`
}
