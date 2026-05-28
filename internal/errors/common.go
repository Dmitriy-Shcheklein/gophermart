// Package errors пакет для доменных ошибок и сообщений
package errors

import "errors"

const (
	// InvalidContentTypeMsg сообщение о невалидном заголовке Content-Type
	InvalidContentTypeMsg = "invalid content-type"
	// DecodeBodyErrMsg сообщение об ошибке декодирования тела запроса
	DecodeBodyErrMsg = "error while decode body"
	// ValidateBodyErrMsg сообщение об ошибке валидации данных
	ValidateBodyErrMsg = "error while validate body"
)

var (
	// ErrEmptyDep ошибка при отсутствии обязательной зависимости
	ErrEmptyDep = errors.New("empty dependency")
	// ErrLoginDuplicate ошибка попытке создать дубликат логина
	ErrLoginDuplicate = errors.New("login already exists")
	// ErrInvalidAuthData ошибка при невалидных данных ждя авторизации
	ErrInvalidAuthData = errors.New("invalid auth data")
	// ErrOrderAlreadyExists ошибка при дубликате заказа
	ErrOrderAlreadyExists = errors.New("order already exists")
	// ErrOrderBelongsAnotherUser ошибка - заказ уже принадлежит другому пользователю
	ErrOrderBelongsAnotherUser = errors.New("order belongs to another user")
	// ErrOrderInvalidNumber ошибка - невалидный номер заказа (не проходит проверку алгоритмом Луна)
	ErrOrderInvalidNumber = errors.New("invalid order number")
	// ErrOrderNotEnoughBalance ошибка - недостаточно средств на балансе
	ErrOrderNotEnoughBalance = errors.New("not enough balance")
	// ErrLoyaltyUnknown неизвестная ошибка от системы расчета
	ErrLoyaltyUnknown = errors.New("loyalty system unknown error")
	// ErrLoyaltyUnknownStatusCode неожиданный статус ответа от системы расчета
	ErrLoyaltyUnknownStatusCode = errors.New("loyalty unknown status code")
	// ErrLoyaltyDecodeBody ошибка при декодировании ответа от системы расчета
	ErrLoyaltyDecodeBody = errors.New("loyalty system decode response error")
	// ErrLoyaltyValidateBody ошибка при валидации ответа от системы расчета
	ErrLoyaltyValidateBody = errors.New("loyalty system validate response error")
	// ErrorLoyaltyWait ожидание тайм-аута для доступа к системе расчета
	ErrorLoyaltyWait = errors.New("loyalty wait error")
	// ErrorLoyaltyTooManyRequest превышен лимит запросов к системе расчета
	ErrorLoyaltyTooManyRequest = errors.New("loyalty too many request")
)
