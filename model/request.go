package model

import (
	"encoding/json"

	"github.com/gitran-com/gitran-server/constant"
)

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
	Token    string   `json:"token"`
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (req *UpdateProfileRequest) Map() map[string]interface{} {
	return map[string]interface{}{
		"name": req.Name,
		"bio":  req.Bio,
	}
}

func (req *UpdateProfileRequest) Valid() bool {
	return req.Name != ""
}

type UpdateProjCfgRequest struct {
	SrcBr         string    `json:"src_br"`
	TrnBr         string    `json:"trn_br"`
	PullGap       uint16    `json:"pull_interval"`
	PushGap       uint16    `json:"push_interval"`
	FileMaps      []FileMap `json:"file_maps"`
	FileMapsBytes []byte    `json:"-"`
	IgnRegs       []string  `json:"ignores"`
	IgnRegsBytes  []byte    `json:"-"`
}

func (req *UpdateProjCfgRequest) Valid() bool {
	var (
		err error
	)
	if req.PullGap < constant.MinProjPullGap &&
		req.PushGap < constant.MinProjPushGap {
		return false
	}
	if req.FileMapsBytes, err = json.Marshal(req.FileMaps); err != nil {
		return false
	}
	if req.IgnRegsBytes, err = json.Marshal(req.IgnRegs); err != nil {
		return false
	}
	return true
}

func (req *UpdateProjCfgRequest) Map() map[string]interface{} {
	json.Marshal(req.FileMaps)
	return map[string]interface{}{
		"src_br":    req.SrcBr,
		"trn_br":    req.TrnBr,
		"pull_gap":  req.PullGap,
		"push_gap":  req.PushGap,
		"file_maps": req.FileMapsBytes,
		"ignores":   req.IgnRegsBytes,
	}
}