package repository

import (
	"fmt"
	e "main/domain/errors"
	"main/domain/model"
	"math/rand"

	"database/sql"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type StoreInterface interface {
	GetUserBannerDB(tagId int, featureId int) (*model.UserBanner, error)
	DeleteBannerDB(id int) error
	FillDB() error
	GetUsers() ([]model.User, error)
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
	 WHERE bt.tag_id = $1 AND bt.feature_id = $2;`, tagId, featureId)
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
	if userBanner.Title == "" {
		return nil, e.ErrNotFound404
	}
	return userBanner, nil
}

func (s *Store) DeleteBannerDB(id int) error {
	_, err := s.db.Exec(`DELETE FROM banners WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ClearDB() error {
	_, err := s.db.Exec("TRUNCATE TABLE users, admins, bannertags, banners, features, tags CASCADE;")
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FillDB() error {
	_ = s.ClearDB()
	tagsCount := 10
	featuresCount := 10
	bannersCount := 10
	for i := 0; i < tagsCount; i++ {
		_, err := s.db.Exec(`INSERT INTO tags (id) VALUES ($1);`, i+1)
		if err != nil {
			_ = s.ClearDB()
			return err
		}
	}
	for i := 0; i < featuresCount; i++ {
		_, err := s.db.Exec(`INSERT INTO features (id) VALUES ($1);`, i+1)
		if err != nil {
			_ = s.ClearDB()
			return err
		}
	}
	_, err := s.db.Exec(`INSERT INTO users (name, token) VALUES ($1, $2);`, "user1", uuid.New().String())
	if err != nil {
		_ = s.ClearDB()
		return err
	}

	_, err = s.db.Exec(`INSERT INTO admins (name, token) VALUES ($1, $2);`, "admin1", uuid.New().String())
	if err != nil {
		_ = s.ClearDB()
		return err
	}
	for i := 0; i < bannersCount; i++ {
		bannerID := 0
		featureID := rand.Intn(featuresCount)
		tagID := rand.Intn(tagsCount)
		err := s.db.QueryRow(`INSERT INTO banners (feature_id, title, text, url) VALUES ($1, $2, $3, $4) RETURNING id;`,
			featureID,
			"title"+fmt.Sprint(i+1),
			"text"+fmt.Sprint(i+1),
			"url"+fmt.Sprint(i+1),
		).Scan(&bannerID)
		if err != nil {
			_ = s.ClearDB()
			return err
		}
		_, err = s.db.Exec(`INSERT INTO bannertags (tag_id, banner_id, feature_id) VALUES ($1, $2, $3);`, tagID, bannerID, featureID)
		if err != nil {
			err2 := s.DeleteBannerDB(bannerID)
			if err2 != nil {
				_ = s.ClearDB()
				return err
			}
		}
	}

	return nil
}

func (s *Store) GetUsers() ([]model.User, error) {
	rows, err := s.db.Query(`SELECT id, name, token FROM users;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var tmpUser model.User
		err := rows.Scan(&tmpUser.Id, &tmpUser.Name, &tmpUser.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, tmpUser)
	}

	rows, err = s.db.Query(`SELECT id, name, token FROM admins;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmpAdmin model.User
		err := rows.Scan(&tmpAdmin.Id, &tmpAdmin.Name, &tmpAdmin.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, tmpAdmin)
	}

	return users, nil
}
