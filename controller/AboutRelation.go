package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"log"
	"songKey/contants"
	"songKey/domain"
	"songKey/services"
)

// CreateRelation need a relationship which need has fromNode and toNode.
func CreateRelation(ctx context.Context, c *app.RequestContext) {
	relationCreate := domain.NewRelation()
	err := c.BindJSON(relationCreate)
	if err != nil {
		log.Println("CreateRelation-Bind Fail")
		c.JSON(contants.ERROR, domain.Response{StatusCode: contants.ERROR_UNMARSHAL_JSON, StatusMsg: err.Error()})
		return
	}
	res, err := services.CreateRelation(relationCreate)
	if err != nil {
		log.Println("CreateRelation-cypher Fail")
		c.JSON(contants.ERROR, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: err.Error()})
	}
	c.JSON(contants.SUCCESS, domain.Response{
		StatusCode: 0,
		StatusMsg:  "success",
		Attach:     utils.H{"result": res},
	})
}
