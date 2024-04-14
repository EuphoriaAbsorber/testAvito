package model

type Error struct {
	Error interface{} `json:"error,omitempty"`
}

type Response struct {
	Body interface{} `json:"body,omitempty"`
}

type UserBanner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}
