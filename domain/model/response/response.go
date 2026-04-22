package response

import "app/domain"

type Base struct {
	ErrorCode  int               `json:"error_code"`
	Message    string            `json:"message"`
	Validation map[string]string `json:"validation,omitempty"`
	Data       any               `json:"data,omitempty"`
}

type List struct {
	List  []any `json:"list"`
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
	Total int64 `json:"total"`
}

// helper functions
func Error(code int, message string) Base {
	return Base{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
	}
}

func ErrorValidation(validation map[string]string, message string) Base {
	return Base{
		ErrorCode:  domain.ErrValidationCode,
		Message:    message,
		Validation: validation,
		Data:       nil,
	}
}

func Success(data any) Base {
	return Base{
		ErrorCode: 0,
		Message:   "OK",
		Data:      data,
	}
}
