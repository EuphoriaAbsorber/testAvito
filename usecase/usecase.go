package usecase

import (
	"main/domain/model"

	rep "main/repository"
)

type UsecaseInterface interface {
	GetUserBanner(tagId int, featureId int, useLastVers bool) (*model.UserBanner, error)
}

type Usecase struct {
	store rep.StoreInterface
}

func NewUsecase(s rep.StoreInterface) UsecaseInterface {
	return &Usecase{
		store: s,
	}
}

func (api *Usecase) GetUserBanner(tagId int, featureId int, useLastVers bool) (*model.UserBanner, error) {
	return api.store.GetUserBannerDB(tagId, featureId)
}
