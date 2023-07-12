package factories

import (
	"context"
	"fmt"

	"github.com/MaheshBailwal/mscore/core"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func CreateServiceContext(ctx *gin.Context) core.ServiceContext {
	sc, _ := ctx.Get("ServiceContext")
	return sc.(core.ServiceContext)
}

func CreateServiceContextGrpc(ctx context.Context) core.ServiceContext {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("meata not found")
	}
	userId := md.Get("user_id")[0]

	return core.ServiceContext{
		CurrentUserId: userId,
		Ctx:           ctx,
	}
}
