package usecase

import (
	"main/domain/model"

	rep "main/repository"
)

type UsecaseInterface interface {
	GetUserBanner(tagId int, featureId int, useLastVers bool) (*model.UserBanner, error)
	FillDB() error
	GetUsers() ([]model.User, error)
}

type Usecase struct {
	store rep.StoreInterface
}

func NewUsecase(s rep.StoreInterface) UsecaseInterface {
	return &Usecase{
		store: s,
	}
}

func (uc *Usecase) GetUserBanner(tagId int, featureId int, useLastVers bool) (*model.UserBanner, error) {
	return uc.store.GetUserBannerDB(tagId, featureId)
}

func (uc *Usecase) FillDB() error {
	return uc.store.FillDB()
}
func (uc *Usecase) GetUsers() ([]model.User, error) {
	return uc.store.GetUsers()
}
