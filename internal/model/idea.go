package model

type IdeaRequest struct {
	Value string `json:"idea"`
}

type IdeaReponse struct {
	Summary     string   `json:"summary"`
	Score       float32  `json:"score"`
	Suggestions []string `json:"suggestions"`
}
