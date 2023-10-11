package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
	"songKey/contants"
	"songKey/domain"
	"songKey/services"
)

func QueryPath(ctx context.Context, c *app.RequestContext) {
	query := domain.NewRelationQuery()
	err := c.BindJSON(query)
	if err != nil {
		log.Println("controller-QueryPath-bindErr:", err)
		c.JSON(consts.StatusAccepted, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "bind error"})
		return
	}
	result, err := services.QueryPath(query)
	if err != nil {
		log.Println("controller-QueryPath-error:", err)
		c.JSON(consts.StatusAccepted, domain.Response{
			StatusCode: contants.ERROR,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": domain.NeoResToResult(result)}})
	return
}
