package repository

import (
	"context"

	"axiata_test/config"
	"axiata_test/model"

	"github.com/google/uuid"
)

type Repository struct {
	DB *config.Database
}

type IRepository interface {
	// posts
	GetAllPosts(ctx context.Context) (*[]model.PostsQueryResponse, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*[]model.PostsQueryResponse, error)
	InsertPost(ctx context.Context, req *model.InsertPostQuery) error
	UpdatePostByID(ctx context.Context, postID uuid.UUID, req *model.InsertPostQuery) error
	DeletePostByID(ctx context.Context, postID uuid.UUID) error
	GetPostByTag(ctx context.Context, tagLabel string) (*[]model.PostsQueryResponse, error)

	// accounts
	RegisterAccount(ctx context.Context, req *model.PayloadRegister) error
	FindUserByUsername(ctx context.Context, username string) (*model.PayloadRegister, error)
}

func NewRepository() IRepository {
	db := config.DatabaseNew()
	return &Repository{
		DB: db,
	}
}