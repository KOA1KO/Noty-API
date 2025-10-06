package pg

import "errors"

var (
	ErrAlredyExists = errors.New("попытка создать существующего пользователя") // 409
	ErrNotFound     = errors.New("запись не найдена")                          // 404
	ErrForbidden    = errors.New("попытка доступа не к своей записи")          // 403
)
