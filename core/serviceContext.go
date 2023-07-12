package core

import "context"

type ServiceContext struct {
	CurrentUserId string
	Ctx           context.Context
	CorrelationId string
}

func NewServiceContext(userId string, ctx context.Context) ServiceContext {
	return ServiceContext{
		CurrentUserId: userId,
		Ctx:           ctx,
	}
}
