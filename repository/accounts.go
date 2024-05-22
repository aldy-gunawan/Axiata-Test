package repository

import (
	"context"

	"axiata_test/model"
)

func (r *Repository) RegisterAccount(ctx context.Context, req *model.PayloadRegister) error {
	// saving the account
	queryPost := `INSERT INTO accounts (username, password, role) VALUES ($1,$2,$3)`
	_, err := r.DB.DB.Exec(queryPost, req.Username, req.Password, req.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*model.PayloadRegister, error) {
	resp := &model.PayloadRegister{}
	var user, pass, rol string
	queryFind := `SELECT username, password, role FROM accounts WHERE username = $1`
	err := r.DB.DB.QueryRow(queryFind, username).Scan(&user, &pass, &rol)
	if err != nil {
		return nil, err
	}
	resp = &model.PayloadRegister{
		Username: user,
		Password: pass,
		Role: rol,
	}
	return resp, nil
}