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

func CreateNode(ctx context.Context, c *app.RequestContext) {
	nodeCreate := domain.NewNode()
	body := make(map[string]interface{})
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(contants.INVALID_PARAMS, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: "no body"})
		return
	}
	lab := body["label"]
	if lab == nil {
		c.JSON(contants.INVALID_PARAMS, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: "no label"})
		return
	}
	nodeCreate.Label = lab.(string)
	propertiesRaw := body["properties"]
	if propertiesRaw == nil {
		c.JSON(contants.INVALID_PARAMS, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: "no properties"})
		return
	}
	properties := propertiesRaw.([]interface{})
	for _, v := range properties {
		nodeCreate.Properties[v.(string)] = "true"
	}
	unique := body["unique"]
	if unique != nil {
		nodeCreate.IsUnique = unique.(string) == "true"
	}
	if nodeCreate.Properties["unique"] == "true" {
		nodeCreate.IsUnique = true
	}
	res, err := services.CreateNode(nodeCreate)
	if err != nil {
		log.Println("CreateNode-cypher Fail")
		c.JSON(contants.ERROR, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: err.Error()})
	}
	c.JSON(contants.SUCCESS, domain.Response{
		StatusCode: 0,
		StatusMsg:  "success",
		Attach:     utils.H{"result": res},
	})
}
