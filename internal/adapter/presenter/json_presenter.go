package presenter

import (
	"encoding/json"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
)

type JSONPresenter struct{}

func NewJSONPresenter() *JSONPresenter { return &JSONPresenter{} }

func (p *JSONPresenter) Present(v any) ([]byte, error) {
	switch t := v.(type) {
	case dto.RegisterOutput:
		return json.Marshal(struct {
			UserID int64  `json:"userId"`
			Name   string `json:"name"`
			Email  string `json:"email"`
		}{UserID: t.UserID, Name: t.Name, Email: t.Email})
	case *dto.RegisterOutput:
		return json.Marshal(struct {
			UserID int64  `json:"userId"`
			Name   string `json:"name"`
			Email  string `json:"email"`
		}{UserID: t.UserID, Name: t.Name, Email: t.Email})
	case dto.LoginOutput:
		return json.Marshal(struct {
			Token string `json:"token"`
		}{Token: t.Token})
	case *dto.LoginOutput:
		return json.Marshal(struct {
			Token string `json:"token"`
		}{Token: t.Token})
	case dto.GetMeOutput:
		return json.Marshal(struct {
			UserID int64  `json:"userId"`
			Name   string `json:"name"`
			Email  string `json:"email"`
		}{UserID: t.UserID, Name: t.Name, Email: t.Email})
	case *dto.GetMeOutput:
		return json.Marshal(struct {
			UserID int64  `json:"userId"`
			Name   string `json:"name"`
			Email  string `json:"email"`
		}{UserID: t.UserID, Name: t.Name, Email: t.Email})
	default:
		return json.Marshal(v)
	}
}
