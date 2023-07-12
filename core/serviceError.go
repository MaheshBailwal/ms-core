package core

import "fmt"

type ServiceError struct {
	Code ErrCode
	Meta any
}

func NewServiceError(code ErrCode, meta any) ServiceError {
	return ServiceError{
		Code: code,
		Meta: meta,
	}
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("Error code %d Meta %s", e.Code, e.Meta)
}
