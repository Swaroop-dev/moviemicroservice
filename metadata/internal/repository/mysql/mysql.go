package mysql

import (
	"database/sql"

	"context"

	_ "github.com/go-sql-driver/mysql"
	"movieapp.com/metadata/internal/repository"
	"movieapp.com/metadata/pkg/model"
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

//method to return the movie records matching the given id

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string

	row := r.db.QueryRowContext(ctx, "SELECT title, description, director FROM movies where id=?", id)

	if err := row.Scan(&title, &description, &director); err != nil {
		return nil, repository.ErrornotFound
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Director:    director,
		Description: description,
	}, nil
}

func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, director, description) VALUES (?,?,?,?)", id, metadata.Title, metadata.Director, metadata.Description)
	if err != nil {
		return err
	}
	return nil
}
