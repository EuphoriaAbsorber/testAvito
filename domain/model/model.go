package model

import "time"

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

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type CreateBanner struct {
	Tag_ids   []int      `json:"tag_ids"`
	FeatureId int        `json:"feature_id"`
	Content   UserBanner `json:"content"`
	IsActive  bool       `json:"is_active"`
}

type Banner struct {
	Id        int        `json:"banner_id"`
	Tag_ids   []int      `json:"tag_ids"`
	FeatureId int        `json:"feature_id"`
	Content   UserBanner `json:"content"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
