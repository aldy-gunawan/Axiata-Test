package usecase

import (
	"context"
	"database/sql"
	// "time"
	
	"axiata_test/model"
	"axiata_test/utils"

	// "github.com/google/uuid"
)

func (uc *Usecase) Login(ctx context.Context, param model.PayloadLogin) model.LoginResponse {
	resp := model.LoginResponse{}
	if param.Username == "" {
		resp.Status.Message = "Username is required!"
		return resp
	}
	if param.Password == "" {
		resp.Status.Message = "Password is required!"
		return resp
	}
	// find accounts by username
	dataUser, err := uc.repository.FindUserByUsername(ctx, param.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			resp.Status.Message = err.Error()
			return resp
		}
		resp.Status.Message = "Username not Found!"
		return resp
	}
	// check the password
	isValid := utils.CheckPasswordHash(param.Password, dataUser.Password)
	if !isValid {
		resp.Status.Message = "Invalid Password"
		return resp
	}
	// generate token jwt
	generatedToken := utils.GenerateTokenJWT(*dataUser)
	resp.Status.Success = true
	resp.Status.Message = "Login Success"
	resp.Token = generatedToken
	return resp
}

func (uc *Usecase) RegisterAccount(ctx context.Context, param model.PayloadRegister) model.CommonResponse {
	resp := model.CommonResponse{}
	if param.Username == "" {
		resp.Message = "Username is required!"
		return resp
	}
	if param.Password == "" {
		resp.Message = "Password is required!"
		return resp
	}
	if param.Role == "" {
		resp.Message = "Role is required!"
		return resp
	}
	if param.Role != "user" && param.Role != "admin" {
		resp.Message = "Role is not available"
		return resp
	}
	// convert data to database requirement
	hashPassword, err := utils.HashPassword(param.Password)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	qryParam := &model.PayloadRegister{
		Username: param.Username,
		Password: hashPassword,
		Role: param.Role,
	}
	err = uc.repository.RegisterAccount(ctx, qryParam)
	if err != nil {
		resp.Message = err.Error()
		return resp
	}
	resp.Success = true
	resp.Message = "success"
	return resp
}