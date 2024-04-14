package usecase

import (
	e "main/domain/errors"
	"main/domain/model"
	"testing"

	mocks "main/mocks"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserBanner(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	StoreMock := mocks.NewMockStoreInterface(ctrl)
	usecase := NewUsecase(StoreMock)

	testUserBanner := new(model.UserBanner)
	err := faker.FakeData(testUserBanner)
	assert.NoError(t, err)

	tagID := 1
	featureID := 1
	useLastVers := true

	//ok
	StoreMock.EXPECT().GetUserBannerDB(tagID, featureID).Return(testUserBanner, nil)
	userBanner, err := usecase.GetUserBanner(tagID, featureID, useLastVers)
	assert.NoError(t, err)
	assert.Equal(t, userBanner, testUserBanner)

	//err
	StoreMock.EXPECT().GetUserBannerDB(tagID, featureID).Return(nil, e.ErrNotFound404)
	_, err = usecase.GetUserBanner(tagID, featureID, useLastVers)
	assert.Equal(t, err, e.ErrNotFound404)
}
