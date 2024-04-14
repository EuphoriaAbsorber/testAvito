package repository

import (
	e "main/domain/errors"
	"main/domain/model"

	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type StoreInterface interface {
	GetUserBannerDB(tagId int, featureId int) (*model.UserBanner, error)
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) StoreInterface {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserBannerDB(tagId int, featureId int) (*model.UserBanner, error) {
	userBanner := &model.UserBanner{}
	rows, err := s.db.Query(`SELECT title, text, url FROM banners b JOIN bannertags bt on b.id = bt.banner_id
	 WHERE bt.tag_id = S1 AND bt.feature_id = $2;`, tagId, featureId)

	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound404
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userBanner.Title, &userBanner.Text, &userBanner.Url)

		if err != nil {
			return nil, err
		}
	}
	return userBanner, nil
}
