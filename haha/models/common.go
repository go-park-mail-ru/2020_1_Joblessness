package models

type ResponseBool struct {
	Like bool `json:"like"`
}

type ResponseID struct {
	ID uint64 `json:"id"`
}
