package usecase

import (
	"context"
	"database/sql"

	"axiata_test/model"

	"github.com/google/uuid"
)

func (uc *Usecase) GetAllTag(ctx context.Context) model.GetAllTagResponse {
	resp := model.GetAllTagResponse{}
	tagsData := []model.TagsData{}
	res, err := uc.repository.GetAllTag(ctx);
	if err != nil {
		if err != sql.ErrNoRows {
			resp.Status.Message = err.Error()
			return resp
		}
		resp.Status.Success = true
		resp.Status.Message = "There are no tags registered yet"
		return resp
	}
	for _, v := range *res {
		eachTag := model.TagsData{
			TagID: v.TagID,
			Label: v.Label,
		}
		tagsData = append(tagsData, eachTag)
	}
	resp.Status.Success = true
	resp.Status.Message = "Success"
	resp.Data = tagsData

	return resp
}

func (uc *Usecase) GetTagByID(ctx context.Context, tagID uuid.UUID) model.GetTagResponse {
	resp := model.GetTagResponse{}

	res, err := uc.repository.GetTagByID(ctx, tagID)
	if err != nil {
		if err != sql.ErrNoRows {
			resp.Status.Message = err.Error()
			return resp
		}
		resp.Status.Message = "Data not found"
		return resp
	}
	resp.Status.Success = true
	resp.Status.Message = "Success"
	resp.Data = *res

	return resp
}

func (uc *Usecase) StoreTag(ctx context.Context, param model.PayloadTags) model.CommonResponse {
	resp := model.CommonResponse{}

	qryParam := &model.PayloadTags{
		Label: param.Label,
	}
	err := uc.repository.InsertTag(ctx, qryParam)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "Success"
	return resp
}

func (uc *Usecase) UpdateTagByID(ctx context.Context, tagID uuid.UUID, param model.PayloadTags) model.CommonResponse {
	resp := model.CommonResponse{}

	qryParam := &model.PayloadTags{
		Label: param.Label,
	}
	err := uc.repository.UpdateTagByID(ctx, tagID, qryParam)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "Success"
	return resp
}

func (uc *Usecase) DeleteTagByID(ctx context.Context, tagID uuid.UUID) model.CommonResponse {
	resp := model.CommonResponse{}
	err := uc.repository.DeleteTagByID(ctx, tagID)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "Success"
	return resp
}