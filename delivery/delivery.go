package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	tok := r.Header.Get("token")
	accessLvl, err := api.usecase.CheckToken(tok)
	if err == e.ErrUnauthorized401 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrUnauthorized401, 401)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
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
	if accessLvl == 1 && (!banner.IsActive) {
		log.Println("banner is unabled for users ")
		ReturnErrorJSON(w, e.ErrForbidden403, 403)
		return
	}

	json.NewEncoder(w).Encode(banner)
}

// GetBanner godoc
// @Summary Получение всех баннеров с фильтрацией по фиче и/или тегу
// @Description Получение всех баннеров с фильтрацией по фиче и/или тегу
// @ID getBanner
// @Accept  json
// @Produce  json
// @Param token header string false "token"
// @Param feature_id query int false "feature_id"
// @Param tag_id query int false "tag_id"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 401 {object} model.Error "Пользователь не авторизован"
// @Failure 403 {object} model.Error "Пользователь не имеет доступа"
// @Failure 404 {object} model.Error "Не найдено"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /banner [get]
func (api *Handler) GetBanners(w http.ResponseWriter, r *http.Request) {
	tagIDs := r.URL.Query().Get("tag_id")
	featureIDs := r.URL.Query().Get("feature_id")
	limitS := r.URL.Query().Get("limit")
	offsetS := r.URL.Query().Get("offset")
	tagID, featureID, limit, offset := 0, 0, 0, 0
	var err error
	if tagIDs != "" {
		tagID, err = strconv.Atoi(tagIDs)
		if err != nil {
			log.Println("error: ", err)
			ReturnErrorJSON(w, e.ErrBadRequest400, 400)
			return
		}
	}
	if featureIDs != "" {
		featureID, err = strconv.Atoi(featureIDs)
		if err != nil {
			log.Println("error: ", err)
			ReturnErrorJSON(w, e.ErrBadRequest400, 400)
			return
		}
	}
	if limitS != "" {
		limit, err = strconv.Atoi(limitS)
		if err != nil {
			log.Println("error: ", err)
			ReturnErrorJSON(w, e.ErrBadRequest400, 400)
			return
		}
	}
	if offsetS != "" {
		offset, err = strconv.Atoi(offsetS)
		if err != nil {
			log.Println("error: ", err)
			ReturnErrorJSON(w, e.ErrBadRequest400, 400)
			return
		}
	}

	tok := r.Header.Get("token")
	accessLvl, err := api.usecase.CheckToken(tok)
	if err == e.ErrUnauthorized401 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrUnauthorized401, 401)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	if accessLvl < 2 {
		log.Println("error: no access")
		ReturnErrorJSON(w, e.ErrForbidden403, 403)
		return
	}

	banners, err := api.usecase.GetBanners(tagID, featureID, limit, offset)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	json.NewEncoder(w).Encode(banners)
}

// CreateBanner godoc
// @Summary Создание нового баннера
// @Description Создание нового баннера
// @ID createBanner
// @Accept  json
// @Produce  json
// @Param banner body model.CreateBanner true "Banner params"
// @Param token header string false "token"
// @Success 201 {object} model.Response "Created"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 401 {object} model.Error "Пользователь не авторизован"
// @Failure 403 {object} model.Error "Пользователь не имеет доступа"
// @Failure 404 {object} model.Error "Не найдено"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /banner [post]
func (api *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.CreateBanner
	err := decoder.Decode(&req)
	if err != nil {
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}

	tok := r.Header.Get("token")
	accessLvl, err := api.usecase.CheckToken(tok)
	if err == e.ErrUnauthorized401 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrUnauthorized401, 401)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	if accessLvl < 2 {
		log.Println("error: no access")
		ReturnErrorJSON(w, e.ErrForbidden403, 403)
		return
	}

	err = api.usecase.CreateBanner(req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	json.NewEncoder(w).Encode(model.Response{})
}

// UpdateBanner godoc
// @Summary Обновление содержимого баннера
// @Description Обновление содержимого баннера
// @ID updateBanner
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param token header string false "token"
// @Param banner body model.CreateBanner true "Banner params"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 401 {object} model.Error "Пользователь не авторизован"
// @Failure 403 {object} model.Error "Пользователь не имеет доступа"
// @Failure 404 {object} model.Error "Не найдено"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /banner/{id}  [patch]
func (api *Handler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	idS := s[len(s)-1]
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req model.CreateBanner
	err = decoder.Decode(&req)
	if err != nil {
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}

	tok := r.Header.Get("token")
	accessLvl, err := api.usecase.CheckToken(tok)
	if err == e.ErrUnauthorized401 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrUnauthorized401, 401)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	if accessLvl < 2 {
		log.Println("error: no access")
		ReturnErrorJSON(w, e.ErrForbidden403, 403)
		return
	}

	err = api.usecase.UpdateBanner(id, req)
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
	json.NewEncoder(w).Encode(model.Response{})
}

// DeleteBanner godoc
// @Summary Удаление баннера по идентификатору
// @Description Удаление баннера по идентификатору
// @ID deleteBanner
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param token header string false "token"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 401 {object} model.Error "Пользователь не авторизован"
// @Failure 403 {object} model.Error "Пользователь не имеет доступа"
// @Failure 404 {object} model.Error "Не найдено"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /banner/{id}  [delete]
func (api *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	idS := s[len(s)-1]
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	tok := r.Header.Get("token")
	accessLvl, err := api.usecase.CheckToken(tok)
	if err == e.ErrUnauthorized401 {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrUnauthorized401, 401)
		return
	}
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrServerError500, 500)
		return
	}
	if accessLvl < 2 {
		log.Println("error: no access")
		ReturnErrorJSON(w, e.ErrForbidden403, 403)
		return
	}

	err = api.usecase.DeleteBanner(id)
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
	json.NewEncoder(w).Encode(model.Response{})
}

// FillDB godoc
// @Summary Заполнение базы тестовыми данными
// @Description Заполнение базы тестовыми данными
// @ID fillDB
// @Accept  json
// @Produce  json
// @Tags extra
// @Param tag_count query int true "tag_count"
// @Param feature_count query int true "feature_count"
// @Param banner_count query int true "banner_count"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Некорректные данные"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /filldb [post]
func (api *Handler) FillDB(w http.ResponseWriter, r *http.Request) {
	tagCountS := r.URL.Query().Get("tag_count")
	featureCountS := r.URL.Query().Get("feature_count")
	bannerCountS := r.URL.Query().Get("banner_count")

	tagCount, err := strconv.Atoi(tagCountS)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	featureCount, err := strconv.Atoi(featureCountS)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	bannerCount, err := strconv.Atoi(bannerCountS)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, e.ErrBadRequest400, 400)
		return
	}
	err = api.usecase.FillDB(tagCount, featureCount, bannerCount)
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
