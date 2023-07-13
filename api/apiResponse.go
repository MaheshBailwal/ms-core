package api

import "github.com/MaheshBailwal/mscore/core"

type ApiResponse[TD any] struct {
	Payload TD `json:"payload"`
} //@name ApiResponse

func newApiResponse[TD any](payload TD) ApiResponse[TD] {
	return ApiResponse[TD]{
		Payload: payload,
	}
}

type ApiError struct {
	ErrorCode core.ErrCode `json:"errorCode" format:"int32"`
	Message   string       `json:"message"`
	Meta      any          `json:"meta,omitempty" extensions:"x-nullable,x-omitempty"`
} //@name ApiError

func NewApiError(code core.ErrCode, message string, meta any) ApiError {
	return ApiError{
		ErrorCode: code,
		Message:   message,
		Meta:      meta,
	}
}

func (r *ApiError) SetMeta(meta any) {
	r.Meta = meta
}
