package usecase

import (
	"strings"
	"context"
	"time"
	
	"axiata_test/model"
	"axiata_test/utils"

	"github.com/google/uuid"
)

func (uc *Usecase) GetAllPost(ctx context.Context) model.GetAllPostsResponse {
	resp := model.GetAllPostsResponse{}
	postsData := []model.PostsData{}
	// get all posts available
	res, err := uc.repository.GetAllPosts(ctx);
	if err != nil {
		resp.Status.Message = err.Error()
		return resp
	}
	mapPost := map[string]model.PostsData{}
	for _, v := range *res {
		if _, ok := mapPost[v.PostID]; !ok {
			tags := []string{}
			postData := model.PostsData{
				PostID: v.PostID,
				Title: v.Title,
				Content: v.Content,
				Status: v.Status,
				PublishDate: v.PublishDate,
				Tags: append(tags, v.Label),
			}
			mapPost[v.PostID] = postData
		} else {
			currentData := mapPost[v.PostID];
			currentTags := currentData.Tags
			currentData.Tags = append(currentTags, v.Label)
			mapPost[v.PostID] = currentData
		}
	}
	for _, mp := range mapPost {
		postsData = append(postsData, mp)
	}
	resp.Status.Success = true
	resp.Status.Message = "Success"
	resp.Data = postsData
	return resp
}

func (uc *Usecase) GetPostByID(ctx context.Context, postID uuid.UUID) model.GetPostsResponse {
	resp := model.GetPostsResponse{}

	res, err := uc.repository.GetPostByID(ctx, postID)
	if err != nil {
		resp.Status.Message = err.Error()
		return resp
	}
	tags := []string{}
	postData := model.PostsData{}
	for _, v := range *res {
		postData = model.PostsData{
			PostID: v.PostID,
			Title: v.Title,
			Content: v.Content,
			Status: v.Status,
			PublishDate: v.PublishDate,
		}
		tags = append(tags, v.Label)
	}
	postData.Tags = tags
	resp.Status.Success = true
	resp.Status.Message = "Success"
	resp.Data = postData
	return resp
}

func (uc *Usecase) StorePosts(ctx context.Context, param model.PostRequest) model.CommonResponse {
	resp := model.CommonResponse{}
	// before inserting, check if the role has authorization
	userRole := ctx.Value(utils.RoleKey).(string)
	userName := ctx.Value(utils.UsernameKey).(string)
	if param.Status && strings.EqualFold(userRole, "user") {
		resp.Message = "Only role ADMIN that can save Post with status Publish"
		return resp
	}
	if !param.Status && strings.EqualFold(userRole, "admin") {
		resp.Message = "Only role USER that can save Post with status Draft"
		return resp
	}

	// convert data to database requirement
	timeNow := time.Now()
	qryParam := &model.InsertPostQuery{
		Title: param.Title,
		Content: param.Content,
		Tags: param.Tags,
		Status: param.Status,
		CreatedBy: userName,
		CreatedDate: timeNow,
	}
	if param.Status {
		qryParam.PublishDate = &timeNow
	}
	err := uc.repository.InsertPost(ctx, qryParam)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "success"
	return resp
}

func (uc *Usecase) UpdatePostByID(ctx context.Context, postID uuid.UUID, param model.PostRequest) model.CommonResponse {
	resp := model.CommonResponse{}
	// before inserting, check if the role has authorization
	userRole := ctx.Value(utils.RoleKey).(string)
	userName := ctx.Value(utils.UsernameKey).(string)
	if param.Status && strings.EqualFold(userRole, "user") {
		resp.Message = "Only role ADMIN that can save Post with status Publish"
		return resp
	}
	if !param.Status && strings.EqualFold(userRole, "admin") {
		resp.Message = "Only role USER that can save Post with status Draft"
		return resp
	}
	timeNow := time.Now()
	qryParam := &model.InsertPostQuery{
		Title: param.Title,
		Content: param.Content,
		Tags: param.Tags,
		Status: param.Status,
		UpdatedBy: &userName,
		UpdatedDate: &timeNow,
	}
	if param.Status && param.PublishDate != nil {
		qryParam.PublishDate = &timeNow
	}
	
	err := uc.repository.UpdatePostByID(ctx, postID, qryParam)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "success"
	return resp
}

func (uc *Usecase) DeletePostByID(ctx context.Context, postID uuid.UUID) model.CommonResponse {
	resp := model.CommonResponse{}
	err := uc.repository.DeletePostByID(ctx, postID)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "Success"
	return resp
}

func (uc *Usecase) GetPostByTag(ctx context.Context, tagLabel string) model.GetAllPostsResponse {
	resp := model.GetAllPostsResponse{}
	postsData := []model.PostsData{}

	res, err := uc.repository.GetPostByTag(ctx, tagLabel)
	if err != nil {
		resp.Status.Message = err.Error()
		return resp
	}
	mapPost := map[string]model.PostsData{}
	for _, v := range *res {
		if _, ok := mapPost[v.PostID]; !ok {
			tags := []string{}
			postData := model.PostsData{
				PostID: v.PostID,
				Title: v.Title,
				Content: v.Content,
				Status: v.Status,
				PublishDate: v.PublishDate,
				Tags: append(tags, v.Label),
			}
			mapPost[v.PostID] = postData
		} else {
			currentData := mapPost[v.PostID];
			currentTags := currentData.Tags
			currentData.Tags = append(currentTags, v.Label)
			mapPost[v.PostID] = currentData
		}
	}
	for _, mp := range mapPost {
		postsData = append(postsData, mp)
	}
	resp.Status.Success = true
	resp.Status.Message = "Success"
	resp.Data = postsData
	return resp
}