package rest

type FileNamesResponse struct {
	Names []string `json:"names"`
}

type FileContentResponse struct {
	Content string `json:"content"`
}
