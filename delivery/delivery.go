package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	e "main/domain/errors"
	"main/domain/model"

	uc "main/usecase"
)

var mockTeacherID = 1

// @title Banner Service API
// @version 1.0
// @description Banner Service back server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath  /api

type Handler struct {
	usecase uc.UsecaseInterface
}

func NewHandler(uc uc.UsecaseInterface) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func ReturnErrorJSON(w http.ResponseWriter, err error, errCode int) {
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(&model.Error{Error: err.Error()})
}

// GetUserBanner godoc
// @Summary Получение баннера для пользователя
// @Description Получение баннера для пользователя
// @ID getUserBanner
// @Accept  json
// @Produce  json
// @Param tag_id query int true "tag_id"
// @Param feature_id query int true "feature_id"
// @Param use_last_revision query boolean false "use_last_revision"
// @Param token header string false "token"
// @Success 200 {object} model.UserBanner "OK"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 401 {object} model.Error "Пользователь не авторизован"
// @Failure 403 {object} model.Error "Пользователь не имеет доступа"
// @Failure 404 {object} model.Error "Не найдено"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /user_banner [get]
func (api *Handler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	tagIDs := r.URL.Query().Get("tag_id")
	featureIDs := r.URL.Query().Get("feature_id")
	useLastRevision := r.URL.Query().Get("use_last_revision")

	tagID, err := strconv.Atoi(tagIDs)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	featureID, err := strconv.Atoi(featureIDs)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	useLastRevisionFlag := false
	if useLastRevision == "true" {
		useLastRevisionFlag = true
	}
	banner, err := api.usecase.GetUserBanner(tagID, featureID, useLastRevisionFlag)
	if err == e.ErrNotFound404 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrNotFound404, 404)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}

	json.NewEncoder(w).Encode(banner)
}

// FillDB godoc
// @Summary Заполнение базы тестовыми данными
// @Description Заполнение базы тестовыми данными
// @ID fillDB
// @Accept  json
// @Produce  json
// @Tags extra
// @Success 200 {object} model.Response "OK"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /filldb [post]
func (api *Handler) FillDB(w http.ResponseWriter, r *http.Request) {

	err := api.usecase.FillDB()
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}

	json.NewEncoder(w).Encode(model.Response{})
}

// GetUsers godoc
// @Summary Получение пользователей и админов
// @Description Получение пользователей и админов
// @ID GetUsers
// @Accept  json
// @Produce  json
// @Tags extra
// @Success 200 {object} model.Response "OK"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /users [get]
func (api *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := api.usecase.GetUsers()
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}

	json.NewEncoder(w).Encode(model.Response{Body: users})
}
