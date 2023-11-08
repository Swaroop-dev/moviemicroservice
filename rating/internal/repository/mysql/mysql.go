package mysql

import (
	"database/sql"

	"context"

	_ "github.com/go-sql-driver/mysql"
	"movieapp.com/rating/internal/repository"
	"movieapp.com/rating/pkg/model"
)

type Repository struct {
	db *sql.DB
}

//New Mysql Based Repo

func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db}, err
}

func (r *Repository) Get(ctx context.Context, recordId model.RecordId, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id,value from ratings where record_id=? AND rating_type=?", recordId, recordType)

	if err != nil {
		return nil, repository.ErrornotFound
	}
	defer rows.Close()
	var res []model.Rating

	for rows.Next() {
		var userId string
		var value int32

		if err := rows.Scan(&userId, &value); err != nil {
			return nil, err
		}
		res = append(res, model.Rating{
			UserId: model.UserId(userId),
			Value:  model.RatingValue(value),
		})
	}

	if len(res) == 0 {
		return nil, repository.ErrornotFound
	}

	return res, nil
}

func (r *Repository) Post(ctx context.Context, recordId model.RecordId, recordtype model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings (record_id,rating_type,user_id,value) VALUES (?,?,?,?,)", recordId, recordtype, rating.UserId, rating.Value)

	if err != nil {
		return err
	}

	return nil
}
