package delivery

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	conf "main/config"
	e "main/domain/errors"
	"main/domain/model"
	mocks "main/mocks"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecaseMock := mocks.NewMockUsecaseInterface(ctrl)
	handler := NewHandler(usecaseMock)

	tagID := 1
	featureID := 1
	useLastVers := true

	testUserBanner := new(model.UserBanner)
	err := faker.FakeData(testUserBanner)
	assert.NoError(t, err)

	tok := "token_string"

	//ok
	url := "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1&use_last_revision=true"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr := httptest.NewRecorder()

	usecaseMock.EXPECT().CheckToken(tok).Return(2, nil)
	usecaseMock.EXPECT().GetUserBanner(tagID, featureID, useLastVers).Return(testUserBanner, nil)
	handler.GetUserBanner(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	//err 400
	url = "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1s&use_last_revision=true"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr = httptest.NewRecorder()
	handler.GetUserBanner(rr, req)
	assert.Equal(t, 400, rr.Code)

	//err 401
	url = "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1&use_last_revision=true"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr = httptest.NewRecorder()
	usecaseMock.EXPECT().CheckToken(tok).Return(0, e.ErrUnauthorized401)
	handler.GetUserBanner(rr, req)
	assert.Equal(t, 401, rr.Code)

	//err 404
	url = "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1&use_last_revision=true"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr = httptest.NewRecorder()
	usecaseMock.EXPECT().CheckToken(tok).Return(2, nil)
	usecaseMock.EXPECT().GetUserBanner(tagID, featureID, useLastVers).Return(nil, e.ErrNotFound404)
	handler.GetUserBanner(rr, req)
	assert.Equal(t, 404, rr.Code)

	//err 500
	url = "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1&use_last_revision=true"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr = httptest.NewRecorder()
	usecaseMock.EXPECT().CheckToken(tok).Return(2, nil)
	usecaseMock.EXPECT().GetUserBanner(tagID, featureID, useLastVers).Return(nil, e.ErrServerError500)
	handler.GetUserBanner(rr, req)
	assert.Equal(t, 500, rr.Code)

	//err 403
	testUserBanner.IsActive = false
	url = "/api" + conf.PathUserBanner + "?tag_id=1&feature_id=1&use_last_revision=true"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", tok)
	rr = httptest.NewRecorder()
	usecaseMock.EXPECT().CheckToken(tok).Return(1, nil)
	usecaseMock.EXPECT().GetUserBanner(tagID, featureID, useLastVers).Return(testUserBanner, nil)
	handler.GetUserBanner(rr, req)
	assert.Equal(t, 403, rr.Code)
}
