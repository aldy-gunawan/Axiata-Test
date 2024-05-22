package repository

import (
	"fmt"
	"time"
	"context"

	"axiata_test/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) GetAllPosts(ctx context.Context) (*[]model.PostsQueryResponse, error) {
	resp := &[]model.PostsQueryResponse{}
	queryGetAll := `SELECT 
		p.id, p.title, p.content, p.status, p.publish_date, t.label 
	from posts p LEFT JOIN posts_tags pt ON p.id = pt.post_id LEFT JOIN tags t ON t.id = pt.tag_id;`
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
			title string
			content string
			status bool
			publish_date *time.Time
			label string
		)
		err = rows.Scan(&id, &title, &content, &status, &publish_date, &label)
		if err != nil {
			return nil, err
		}
		*resp = append(*resp, model.PostsQueryResponse{
			PostID: id,
			Title: title,
			Content: content,
			Status: status,
			PublishDate: publish_date,
			Label: label,
		})
	}
	defer rows.Close()

	return resp, nil
}

func (r *Repository) GetPostByID(ctx context.Context, postID uuid.UUID) (*[]model.PostsQueryResponse, error) {
	resp := &[]model.PostsQueryResponse{}
	queryGetByID := `SELECT 
		p.id, p.title, p.content, p.status, p.publish_date, t.label 
	from posts p LEFT JOIN posts_tags pt ON p.id = pt.post_id LEFT JOIN tags t ON t.id = pt.tag_id
	WHERE p.id = $1`
	rows, err := r.DB.DB.QueryContext(ctx, queryGetByID, postID)
	if err != nil {
		return nil, err
	}
    if rows == nil {
        return nil, nil
    }
	for rows.Next() {
		var (
			id string
			title string
			content string
			status bool
			publish_date *time.Time
			label string
		)
		err = rows.Scan(&id, &title, &content, &status, &publish_date, &label)
		if err != nil {
			return nil, err
		}
		*resp = append(*resp, model.PostsQueryResponse{
			PostID: id,
			Title: title,
			Content: content,
			Status: status,
			PublishDate: publish_date,
			Label: label,
		})
	}
	defer rows.Close()

	return resp, nil
}

func (r *Repository) InsertPost(ctx context.Context, req *model.InsertPostQuery) error {
	// Saving the post
	postID := uuid.New()
	queryPost := `INSERT INTO posts (id, title, content, status, publish_date, created_by, created_date)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.DB.DB.Exec(queryPost, postID, req.Title, req.Content, req.Status, req.PublishDate, req.CreatedBy, req.CreatedDate)
	if err != nil {
		return err
	}

	// saving the tags
	for _, v := range req.Tags {
		newID := uuid.New()
		queryTag := `INSERT INTO tags (id, label) VALUES ($1,$2) ON CONFLICT ON CONSTRAINT label_uniq DO NOTHING`
		_, err := r.DB.DB.Exec(queryTag, newID, v)
		if err != nil {
			return err
		}
	}

	// finding the tags id after inserting it
	queryFindTag := `SELECT id FROM tags WHERE label = ANY($1)`
	rows, err := r.DB.DB.QueryContext(ctx, queryFindTag, pq.Array(req.Tags))
	if err != nil {
		return err
	}
    if rows == nil {
        fmt.Println("no rows found")
    }
	var tagIDs []uuid.UUID
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, uuid.MustParse(id))
	}
	defer rows.Close()

	// saving the mapping of posts and tags
	for _, v := range tagIDs {
		queryPostTag := `INSERT INTO posts_tags (post_id, tag_id) VALUES ($1,$2)`
		_, err = r.DB.DB.Exec(queryPostTag, postID, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdatePostByID(ctx context.Context, postID uuid.UUID, req *model.InsertPostQuery) error {
	// begin transaction
	tx, err := r.DB.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
        if err != nil {
            fmt.Println("Rolling back transaction due to error:", err)
            if rollbackErr := tx.Rollback(); rollbackErr != nil {
                fmt.Println("Error rolling back transaction:", rollbackErr)
            }
        } else {
            fmt.Println("Committing transaction")
            if commitErr := tx.Commit(); commitErr != nil {
                fmt.Println("Error committing transaction:", commitErr)
            }
        }
    }()

	// delete all relations posts_tags
	queryDelete := `DELETE FROM posts_tags WHERE post_id = $1`
	_, err = tx.Exec(queryDelete, postID)
	if err != nil {
		return err
	}

	// update posts
	queryUpdate := `UPDATE posts SET title = $2, content = $3, status = $4, publish_date = $5, updated_date = $6, updated_by = $7 WHERE id = $1`
	_, err = tx.Exec(queryUpdate, postID, req.Title, req.Content, req.Status, req.PublishDate, req.UpdatedDate, req.UpdatedBy)
	if err != nil {
		return err
	}

	// inserting if there are any new tags
	tagsID := []uuid.UUID{}
	for _, v := range req.Tags {
		newID := uuid.New()
		var stringID string
		queryTag := `INSERT INTO tags (id, label) VALUES ($1,$2) ON CONFLICT ON CONSTRAINT label_uniq DO UPDATE SET label = $2 RETURNING id`
		err = tx.QueryRow(queryTag, newID, v).Scan(&stringID)
		if err != nil {
			return err
		}
		tagsID = append(tagsID, uuid.MustParse(stringID))
	}

	// saving the mapping of posts and tags
	for _, v := range tagsID {
		queryPostTag := `INSERT INTO posts_tags (post_id, tag_id) VALUES ($1,$2)`
		_, err = tx.Exec(queryPostTag, postID, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeletePostByID(ctx context.Context, postID uuid.UUID) error {
	queryDelete := `DELETE FROM posts WHERE id = $1`
	_, err := r.DB.DB.Exec(queryDelete, postID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPostByTag(ctx context.Context, tagLabel string) (*[]model.PostsQueryResponse, error) {
	resp := &[]model.PostsQueryResponse{}
	queryGetByTagLabel := `
		SELECT po.id, po.title, po.content, po.status, po.publish_date, ta.label
		FROM posts po
		LEFT JOIN posts_tags pot ON po.id = pot.post_id
		LEFT JOIN tags ta ON ta.id = pot.tag_id
		WHERE po.id = ANY(SELECT pt.post_id AS id FROM posts_tags pt
			LEFT JOIN tags t ON pt.tag_id = t.id
			WHERE t.label = $1)
		ORDER BY po.created_date ASC;
	`
	rows, err := r.DB.DB.QueryContext(ctx, queryGetByTagLabel, tagLabel)
	if err != nil {
		return nil, err
	}
    if rows == nil {
        return nil, nil
    }
	for rows.Next() {
		var (
			id string
			title string
			content string
			status bool
			publish_date *time.Time
			label string
		)
		err = rows.Scan(&id, &title, &content, &status, &publish_date, &label)
		if err != nil {
			return nil, err
		}
		*resp = append(*resp, model.PostsQueryResponse{
			PostID: id,
			Title: title,
			Content: content,
			Status: status,
			PublishDate: publish_date,
			Label: label,
		})
	}
	defer rows.Close()

	return resp, nil
}