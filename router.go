package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"songKey/controller"
	"songKey/myUtils"
)

func StartRouter() {
	r := server.Default()
	r.Use(myUtils.GlobalErrorHandler)

	router(r)
	r.Spin()
}

func router(r *server.Hertz) {
	r.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
	r.POST("/createRelation", controller.CreateRelation)
	r.POST("/createNode", controller.CreateNode)
}
