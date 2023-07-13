package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MaheshBailwal/mscore/core"
	"github.com/MaheshBailwal/mscore/persistence"
	"github.com/ehsandavari/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController struct {
}

func NewBaseController() BaseController {
	return BaseController{}
}

type RequestHandler[TReq, TRes any] func(ctx *gin.Context, request TReq) (TRes, error)

func (r RequestHandler[TReq, TRes]) Handle(iLogger logger.ILogger, errorMsgProvider IErrorMessageProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TReq

		if bindErr := ctx.ShouldBind(&request); bindErr != nil {
			iLogger.Error(ctx, bindErr.Error())
			err := NewApiError(http.StatusBadRequest, "error in validate request", nil)
			err.SetMeta(bindErr.Error())
			if validationErrors, ok := bindErr.(validator.ValidationErrors); ok {
				meta := make(map[string]string, len(validationErrors))
				for _, validationError := range validationErrors {
					meta[validationError.Field()] = validationError.Field() + " is " + validationError.Tag() + " " + validationError.Param()
				}
				err.SetMeta(meta)
			}
			ctx.JSON(http.StatusBadRequest, newApiResponse[ApiError](
				err,
			))
			return
		}

		result, err := r(ctx, request)

		if err != nil {
			switch err.(type) {
			case core.ServiceError:
				if e, ok := err.(core.ServiceError); ok {
					status := errorMsgProvider.GetHttpStatus(e)
					ctx.JSON(status, newApiResponse[ApiError](
						NewApiError(e.Code, errorMsgProvider.GetErrorMessage()(e.Code), e.Meta),
					))
				}
			}
			return
		}

		ctx.JSON(http.StatusOK, newApiResponse[TRes](
			result,
		))
	}
}

func (r *BaseController) ParseFilter(filterQuery string) persistence.QueryFilter {
	queryFilter := persistence.QueryFilter{}
	queryFilter.Filters = make(map[string]string)
	queryFilter.Sort = make(map[string]string)
	if filterQuery == "" {
		return queryFilter
	}

	for _, filter := range strings.Split(filterQuery, ",") {
		arr := strings.Split(filter, "=")
		fmt.Println(arr)
		if strings.ToLower(arr[0]) == "sortby" {
			queryFilter.Sort[strings.ToLower(arr[1])] = "ASC"
		}
		queryFilter.Filters[strings.ToLower(arr[0])] = arr[1]
	}

	fmt.Println(queryFilter)
	return queryFilter
}
