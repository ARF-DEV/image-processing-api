package model

type Meta struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}
