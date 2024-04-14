package repository

import (
	"main/domain/model"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserBannerDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var tagID int = 1
	var featureID int = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"title", "text", "url", "is_active"})
	expect := &model.UserBanner{Title: "title1", Text: "text1", Url: "url1", IsActive: true}

	rows = rows.AddRow(expect.Title, expect.Text, expect.Url, expect.IsActive)

	mock.
		ExpectQuery(`SELECT title, text, url, is_active FROM banners b JOIN bannertags bt on b.id = bt.banner_id WHERE`).
		WithArgs(tagID, featureID).
		WillReturnRows(rows)

	repo := &Store{
		db: db,
	}
	item, err := repo.GetUserBannerDB(tagID, featureID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect) {
		t.Errorf("results not match, want %v, have %v", expect, item)
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"title", "text"}).
		AddRow("title2", "text2")

	mock.
		ExpectQuery(`SELECT title, text, url, is_active FROM banners b JOIN bannertags bt on b.id = bt.banner_id WHERE`).
		WithArgs(tagID, featureID).
		WillReturnRows(rows)

	_, err = repo.GetUserBannerDB(tagID, featureID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// not found error
	rows = sqlmock.NewRows([]string{"title", "text", "url", "is_active"})
	mock.
		ExpectQuery(`SELECT title, text, url, is_active FROM banners b JOIN bannertags bt on b.id = bt.banner_id WHERE`).
		WithArgs(tagID, featureID).
		WillReturnRows(rows)

	_, err = repo.GetUserBannerDB(tagID, featureID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
