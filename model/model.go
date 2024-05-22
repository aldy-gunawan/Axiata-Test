package model

import (
	"time"
)

type CommonResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}


type GetAllPostsResponse struct {
	Status CommonResponse `json:"status"`
	Data []PostsData `json:"data"`
}

type GetPostsResponse struct {
	Status CommonResponse `json:"status"`
	Data PostsData `json:"data"`
}

type PostsData struct {
	PostID string `json:"post_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status bool `json:"status"`
	Tags []string `json:"tags"`
	PublishDate *time.Time `json:"publish_date"`
}

type PostRequest struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Tags []string `json:"tags"`
	Status bool `json:"status"`
	PublishDate *time.Time `json:"publish_date"`
}

type InsertPostQuery struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Status bool `json:"status"`
	Tags []string `json:"tags"`
	PublishDate *time.Time `json:"publish_date"`
	CreatedDate time.Time `json:"created_date"`
	CreatedBy string `json:"created_by"`
	UpdatedDate *time.Time `json:"updated_date"`
	UpdatedBy *string `json:"updated_by"`
}

type PostsQueryResponse struct {
	PostID string `json:"post_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status bool `json:"status"`
	Label string `json:"label"`
	PublishDate *time.Time `json:"publish_date"`
}

type PayloadRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type PayloadLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status CommonResponse `json:"status"`
	Token string `json:"token"`
}