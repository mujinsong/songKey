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

func SetNode(ctx context.Context, c *app.RequestContext) {
	node := domain.NewNode()
	err := c.BindJSON(node)
	if err != nil {
		log.Println("SetNode-bind-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，bind解析失败"})
		return
	}
	if node.Id == -1 {
		log.Println("SetNode: dont get the id of node because node.Id == -1")
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，没传id，id为-1"})
		return
	}
	result, err := services.SetNode(node)
	if err != nil {
		log.Println("setNode-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.ERROR, StatusMsg: "setNode失败"})
		return
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": result}})
}
