package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("неверные логин/пароль")      // 409
	ErrUserExists         = errors.New("такое имя уже используется") // 409
	ErrUserNotFound       = errors.New("пользователь не найдет")     // 404
	ErrForbidden          = errors.New("доступ закрыт")              // 403

	ErrNoteNotFound = errors.New("заметка не найдена") // 404

	ErrBadRequest = errors.New("некорректный запрос") // 400
	ErrInternal   = errors.New("внутренняя ошибка")   // 500
)
