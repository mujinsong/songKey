package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"songKey/contants"
	"songKey/domain"
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
	cy := domain.CypherStruct{}
	cypher := cy.MatchNode(relationCreate.FromNode).MatchNode(relationCreate.ToNode).WhereAnd("n0", relationCreate.FromNode).WhereAnd("n1", relationCreate.ToNode).CreateOnlyRelation(relationCreate, "n0", "n1").ReturnAll().GetFinalCypher()
	log.Println("CreateRelation:" + cypher)
	res, err := cy.Result()
	if err != nil {
		log.Println("CreateRelation-cypher Fail")
		c.JSON(contants.ERROR, domain.Response{StatusCode: contants.ERROR_CALL_API, StatusMsg: err.Error()})
	}
	c.JSON(contants.SUCCESS, domain.Response{
		StatusCode: 0,
		StatusMsg:  "success",
		Attach:     res,
	})
}
