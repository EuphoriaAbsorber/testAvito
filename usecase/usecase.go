package usecase

import (
	"main/domain/model"

	rep "main/repository"
)

type UsecaseInterface interface {
	GetUserBanner(tagId int, featureId int, useLastVers bool) (*model.UserBanner, error)
	GetBanners(tagId int, featureId int, limit int, offset int) ([]model.Banner, error)
	CreateBanner(req model.CreateBanner) error
	UpdateBanner(id int, req model.CreateBanner) error
	DeleteBanner(id int) error
	CheckToken(tok string) (int, error)
	FillDB(tagCount int, featureCount int, bannerCount int) error
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
func (uc *Usecase) GetBanners(tagId int, featureId int, limit int, offset int) ([]model.Banner, error) {
	return uc.store.GetBannersDB(tagId, featureId, limit, offset)
}
func (uc *Usecase) CreateBanner(req model.CreateBanner) error {
	return uc.store.CreateBannerDB(req)
}
func (uc *Usecase) UpdateBanner(id int, req model.CreateBanner) error {
	return uc.store.UpdateBannerDB(id, req)
}
func (uc *Usecase) DeleteBanner(id int) error {
	return uc.store.DeleteBannerDB(id)
}

func (uc *Usecase) CheckToken(tok string) (int, error) {
	return uc.store.CheckTokenDB(tok)
}

func (uc *Usecase) FillDB(tagCount int, featureCount int, bannerCount int) error {
	return uc.store.FillDB(tagCount, featureCount, bannerCount)
}
func (uc *Usecase) GetUsers() ([]model.User, error) {
	return uc.store.GetUsersDB()
}
