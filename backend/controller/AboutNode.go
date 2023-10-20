package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
	"songKey/contants"
	"songKey/domain"
	"songKey/services"
)

func CreateNode(ctx context.Context, c *app.RequestContext) {
	temp := make([]domain.Node, 1)
	body, _ := c.Body()
	err := json.Unmarshal(body, &temp)
	if err != nil {
		log.Println("CreateNode-unmarshal-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，解析失败"})
		return
	}
	l := len(temp)
	nodes := make([]*domain.Node, l)
	for i := 0; i < l; i++ {
		nodes[i] = domain.NewNode()
	}
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Println("CreateNode-unmarshal-err:", err)
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

func SetNode(ctx context.Context, c *app.RequestContext) {
	temp := make([]domain.Node, 1)
	body, _ := c.Body()
	err := json.Unmarshal(body, &temp)
	if err != nil {
		log.Println("SetNode-unmarshal-err:", err)
		c.JSON(contants.SUCCESS, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，解析失败"})
		return
	}
	l := len(temp)
	nodes := make([]*domain.Node, l)
	for i := 0; i < l; i++ {
		nodes[i] = domain.NewNode()
	}
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Println("SetNode-unmarshal-err:", err)
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
