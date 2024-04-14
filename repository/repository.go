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
	GetBannersDB(tagId int, featureId int, limit int, offset int) ([]model.Banner, error)
	CreateBannerDB(req model.CreateBanner) error
	UpdateBannerDB(id int, req model.CreateBanner) error
	DeleteBannerDB(id int) error
	FillDB(tagCount int, featureCount int, bannerCount int) error
	GetUsersDB() ([]model.User, error)
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

func (s *Store) GetBannersDB(tagId int, featureId int, limit int, offset int) ([]model.Banner, error) {
	if limit == 0 {
		limit = 10000000
	}
	neededBannerIDs := []int{}
	banners := []model.Banner{}
	var rows *sql.Rows
	var err error = nil
	tmpId := 0
	//defer rows.Close()
	if tagId != 0 {
		if featureId != 0 {
			rows, err = s.db.Query(`SELECT banner_id FROM bannertags WHERE tag_id = $1 AND feature_id = $2;`, tagId, featureId)
		} else if featureId == 0 {
			rows, err = s.db.Query(`SELECT banner_id FROM bannertags WHERE tag_id = $1;`, tagId)
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmpId)
			if err != nil {
				return nil, err
			}
			neededBannerIDs = append(neededBannerIDs, tmpId)
		}
		counter := 0
		for ind := range neededBannerIDs {
			counter++
			if ind+offset >= len(neededBannerIDs) {
				break
			}
			tmpBanner := model.Banner{}
			err = s.db.QueryRow(`SELECT id, feature_id, title, text, url, is_active, created_at, updated_at FROM banners WHERE id = $1;`,
				neededBannerIDs[ind+offset]).Scan(&tmpBanner.Id, &tmpBanner.FeatureId, &tmpBanner.Content.Title, &tmpBanner.Content.Text, &tmpBanner.Content.Url,
				&tmpBanner.IsActive, &tmpBanner.CreatedAt, &tmpBanner.UpdatedAt)
			banners = append(banners, tmpBanner)
			if counter >= limit {
				break
			}
		}
		if err != nil {
			return nil, err
		}

	} else {
		if featureId != 0 {
			rows, err = s.db.Query(`SELECT id, feature_id, title, text, url, is_active, created_at, updated_at FROM banners WHERE feature_id = $1 LIMIT $2 OFFSET $3;`, featureId, limit, offset)
		} else if featureId == 0 {
			rows, err = s.db.Query(`SELECT id, feature_id, title, text, url, is_active, created_at, updated_at FROM banners LIMIT $1 OFFSET $2;`, limit, offset)
		}
		defer rows.Close()
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			tmpBanner := model.Banner{}
			err := rows.Scan(&tmpBanner.Id, &tmpBanner.FeatureId, &tmpBanner.Content.Title, &tmpBanner.Content.Text, &tmpBanner.Content.Url,
				&tmpBanner.IsActive, &tmpBanner.CreatedAt, &tmpBanner.UpdatedAt)
			if err != nil {
				return nil, err
			}
			banners = append(banners, tmpBanner)
		}
	}
	for ind, banner := range banners {
		rows, err = s.db.Query(`SELECT tag_id FROM bannertags WHERE banner_id = $1 AND feature_id = $2;`, banner.Id, banner.FeatureId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmpId)
			if err != nil {
				return nil, err
			}
			banners[ind].Tag_ids = append(banners[ind].Tag_ids, tmpId)
		}
	}
	return banners, nil
}

func (s *Store) CreateBannerDB(req model.CreateBanner) error {
	bannerID := 0
	err := s.db.QueryRow(`INSERT INTO banners (feature_id, title, text, url, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		req.FeatureId,
		req.Content.Title,
		req.Content.Text,
		req.Content.Url,
		req.IsActive,
	).Scan(&bannerID)
	if err != nil {
		return err
	}
	for _, tagID := range req.Tag_ids {
		_, err = s.db.Exec(`INSERT INTO bannertags (tag_id, banner_id, feature_id) VALUES ($1, $2, $3);`, tagID, bannerID, req.FeatureId)
		if err != nil {
			_ = s.DeleteBannerDB(bannerID)
			return err
		}
	}
	return nil
}

func (s *Store) UpdateBannerDB(id int, req model.CreateBanner) error {
	bannerID := 0
	err := s.db.QueryRow(`SELECT id FROM banners WHERE id = $1;`, id).Scan(&bannerID)
	if err == sql.ErrNoRows {
		return e.ErrNotFound404
	}
	if err != nil {
		return err
	}
	if bannerID != id {
		return e.ErrNotFound404
	}
	_, err = s.db.Exec(`UPDATE banners SET feature_id = $1, title = $2, text = $3, url = $4, is_active = $5, updated_at = now() WHERE id = $6;`,
		req.FeatureId,
		req.Content.Title,
		req.Content.Text,
		req.Content.Url,
		req.IsActive,
		id,
	)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM bannertags WHERE banner_id = $1;`, bannerID)
	if err != nil {
		return err
	}

	for _, tagID := range req.Tag_ids {
		_, err = s.db.Exec(`INSERT INTO bannertags (tag_id, banner_id, feature_id) VALUES ($1, $2, $3);`, tagID, bannerID, req.FeatureId)
		if err != nil {
			_ = s.DeleteBannerDB(bannerID)
			return err
		}
	}
	return nil
}

func (s *Store) DeleteBannerDB(id int) error {
	bannerID := 0
	err := s.db.QueryRow(`SELECT id FROM banners WHERE id = $1;`, id).Scan(&bannerID)
	if err == sql.ErrNoRows {
		return e.ErrNotFound404
	}
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM banners WHERE id = $1;`, id)
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

func (s *Store) FillDB(tagCount int, featureCount int, bannerCount int) error {
	_ = s.ClearDB()
	//tagsCount := 10
	//featuresCount := 10
	//bannersCount := 10
	for i := 0; i < tagCount; i++ {
		_, err := s.db.Exec(`INSERT INTO tags (id) VALUES ($1);`, i+1)
		if err != nil {
			_ = s.ClearDB()
			return err
		}
	}
	for i := 0; i < featureCount; i++ {
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
	for i := 0; i < bannerCount; i++ {
		bannerID := 0
		featureID := rand.Intn(featureCount)
		tagID := rand.Intn(tagCount)
		err := s.db.QueryRow(`INSERT INTO banners (feature_id, title, text, url) VALUES ($1, $2, $3, $4) RETURNING id;`,
			featureID,
			"title"+fmt.Sprint(i+1),
			"text"+fmt.Sprint(i+1),
			"url"+fmt.Sprint(i+1),
		).Scan(&bannerID)
		if err != nil {
			continue
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

func (s *Store) GetUsersDB() ([]model.User, error) {
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
