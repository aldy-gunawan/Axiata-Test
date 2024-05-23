package repository

import (
	"context"

	"axiata_test/model"

	"github.com/google/uuid"
)

func (r *Repository) GetAllTag(ctx context.Context) (*[]model.TagsData, error) {
	resp := &[]model.TagsData{}
	queryGetAll := `SELECT id, label FROM tags`
	rows, err := r.DB.DB.QueryContext(ctx, queryGetAll)
	if err != nil {
		return nil, err
	}
    if rows == nil {
        return nil, nil
    }
	for rows.Next() {
		var (
			id string
			label string
		)
		err = rows.Scan(&id, &label)
		if err != nil {
			return nil, err
		}
		*resp = append(*resp, model.TagsData{
			TagID: id,
			Label: label,
		})
	}
	defer rows.Close()

	return resp, nil
}

func (r *Repository) GetTagByID(ctx context.Context, tagID uuid.UUID) (*model.TagsData, error) {
	resp := &model.TagsData{}
	queryFind := `SELECT id, label FROM tags WHERE id = $1 `
	var tagId, label string
	err := r.DB.DB.QueryRow(queryFind, tagID).Scan(&tagId, &label)
	if err != nil {
		return nil, err
	}
	resp = &model.TagsData{
		TagID: tagId,
		Label: label,
	}

	return resp, nil
}

func (r *Repository) InsertTag(ctx context.Context, req *model.PayloadTags) error {
	// saving tag
	tagID := uuid.New()
	queryPost := `INSERT INTO tags (id, label) VALUES ($1,$2)`
	_, err := r.DB.DB.Exec(queryPost, tagID, req.Label)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateTagByID(ctx context.Context, tagID uuid.UUID, req *model.PayloadTags) error {
	// update tag
	queryUpdate := `UPDATE tags SET label = $2 WHERE id = $1`
	_, err := r.DB.DB.Exec(queryUpdate, tagID, req.Label)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteTagByID(ctx context.Context, tagID uuid.UUID) error {
	queryDelete := `DELETE FROM tags WHERE id = $1`
	_, err := r.DB.DB.Exec(queryDelete, tagID)
	if err != nil {
		return err
	}
	return nil
}