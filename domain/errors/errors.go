package Errors

import (
	"errors"
)

var ErrBadRequest400 = errors.New("Некорректные данные")
var ErrUnauthorized401 = errors.New("Пользователь не авторизован")
var ErrForbidden403 = errors.New("Пользователь не имеет доступа")
var ErrNotFound404 = errors.New("Не найдено")
var ErrServerError500 = errors.New("Внутренняя ошибка сервера")
