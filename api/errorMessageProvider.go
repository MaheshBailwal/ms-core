package api

import "github.com/MaheshBailwal/mscore/core"

type IErrorMessageProvider interface {
	GetErrorMessage() func(core.ErrCode) string
	GetHttpStatus(err core.ServiceError) int
}
