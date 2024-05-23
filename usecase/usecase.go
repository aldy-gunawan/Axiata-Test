package usecase

import (
	"context"

	"axiata_test/repository"
	"axiata_test/model"
	
	"github.com/google/uuid"
)

type Usecase struct {
	repository repository.IRepository
}

type IUsecase interface {
	// posts
	GetAllPost(ctx context.Context) model.GetAllPostsResponse
	GetPostByID(ctx context.Context, postID uuid.UUID) model.GetPostsResponse
	StorePosts(ctx context.Context, param model.PostRequest) model.CommonResponse
	UpdatePostByID(ctx context.Context, postID uuid.UUID, param model.PostRequest) model.CommonResponse
	DeletePostByID(ctx context.Context, postID uuid.UUID) model.CommonResponse
	GetPostByTag(ctx context.Context, tagLabel string) model.GetAllPostsResponse

	// accounts
	RegisterAccount(ctx context.Context, param model.PayloadRegister) model.CommonResponse
	Login(ctx context.Context, param model.PayloadLogin) model.LoginResponse

	// tags
	GetAllTag(ctx context.Context) model.GetAllTagResponse
	GetTagByID(ctx context.Context, tagID uuid.UUID) model.GetTagResponse
	StoreTag(ctx context.Context, param model.PayloadTags) model.CommonResponse
	UpdateTagByID(ctx context.Context, tagID uuid.UUID, param model.PayloadTags) model.CommonResponse
	DeleteTagByID(ctx context.Context, tagID uuid.UUID) model.CommonResponse
}

func New() IUsecase {
	return &Usecase{
		repository: repository.NewRepository(),
	}
}
