package dto

type PageDTO[T any] struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
	Data   []T `json:"data"`
}
