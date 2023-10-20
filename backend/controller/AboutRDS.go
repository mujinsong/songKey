package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
	"songKey/contants"
	"songKey/dao/rds"
	"songKey/domain"
	"songKey/global"
	"songKey/services"
)

func GetTableMes(ctx context.Context, c *app.RequestContext) {
	tableMes := rds.GetTableMes()
	c.JSON(consts.StatusOK, domain.Response{
		StatusCode: consts.StatusOK,
		StatusMsg:  "done",
		Attach:     utils.H{"result": tableMes},
	})
}

func GetFieldMes(ctx context.Context, c *app.RequestContext) {
	tables := make([]string, 1)
	err := c.BindJSON(&tables)
	if err != nil {
		log.Println("GetFieldMes-bind-err:", err)
		c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.INVALID_PARAMS, StatusMsg: "传参错误，bind解析失败"})
		return
	}
	result := services.GetFieldMes(tables)
	c.JSON(consts.StatusOK, domain.Response{
		StatusCode: consts.StatusOK,
		StatusMsg:  "done",
		Attach:     utils.H{"result": result},
	})
}

func ChangeRdsDb(ctx context.Context, c *app.RequestContext) {
	name := c.Query("dbName")
	if name == "" {
		log.Println("get param error")
		c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.ERROR, StatusMsg: "dbName is nil"})
		return
	}
	log.Println("Now DB is", global.RdsDbName, "wanna change to", name)
	db, err := rds.ChangeDb(name)
	if err != nil {
		log.Println("change dbName error")
		c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.ERROR, StatusMsg: "传参错误"})
		return
	}
	if !db {
		log.Println("change dbName fail")
		c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.ERROR, StatusMsg: "fail"})
		return
	}
	log.Println("changed DB to", global.RdsDbName)
	c.JSON(consts.StatusOK, domain.Response{StatusCode: contants.SUCCESS, StatusMsg: "ok", Attach: utils.H{"result": global.RdsDbName}})
}
