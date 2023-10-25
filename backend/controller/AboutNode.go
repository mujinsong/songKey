package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
	"songKey/contants"
	"songKey/domain"
	"songKey/myUtils"
	"songKey/services"
)

func CreateNode(ctx context.Context, c *app.RequestContext) {
	body, _ := c.Body()
	nodes, err := myUtils.NodesGet(body)
	if err != nil {
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，解析失败"})
		return
	}
	result := make([]*domain.Result, 0)
	for _, node := range nodes {
		res, err := services.CreateNode(node)
		if err != nil {
			log.Println("CreateNode-err:", err)
			continue
		}
		result = append(result, res)
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": result}})
}

func SetNodes(ctx context.Context, c *app.RequestContext) {
	body, _ := c.Body()
	nodes, err := myUtils.NodesGet(body)
	if err != nil {
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，解析失败"})
		return
	}
	result := make([]*domain.Result, 0)
	for _, node := range nodes {
		if node.Id == -1 {
			log.Println("SetNode: dont get the id of node because node.Id == -1")
			continue
		}
		res, err := services.SetNode(node)
		if err != nil {
			log.Println("setNode-err:", err)
			continue
		}
		result = append(result, res)
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": result}})
}
func MatchNode(ctx context.Context, c *app.RequestContext) {
	node := domain.NewNode()
	err := c.BindJSON(node)
	if err != nil {
		log.Println("MatchNode-bind-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，bind解析失败"})
		return
	}
	result, err := services.MatchNode(node)
	if err != nil {
		log.Println("setNode-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.ERROR, StatusMsg: "setNode失败"})
		return
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": result}})
}

func DeleteNodes(ctx context.Context, c *app.RequestContext) {
	body, _ := c.Body()
	nodes, err := myUtils.NodesGet(body)
	if err != nil {
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，解析失败"})
		return
	}
	result, err := services.DeleteNode(nodes)
	if err != nil {
		log.Println("deleteNodes-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.ERROR, StatusMsg: "deleteNodes失败"})
		return
	}
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, Attach: utils.H{"result": result}})
}
