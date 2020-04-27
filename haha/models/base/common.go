package baseModels

//easyjson:json
type ResponseBool struct {
	Like bool `json:"like"`
}

//easyjson:json
type ResponseID struct {
	ID uint64 `json:"id"`
}
