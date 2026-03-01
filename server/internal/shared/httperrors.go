package shared

type AppError struct {
	Status int
	Code   ErrorCode
	Msg    string
}

func (e *AppError) Error() string { return e.Msg }

func NewAppError(status int, code ErrorCode, msg string) *AppError {
	return &AppError{Status: status, Code: code, Msg: msg}
}
