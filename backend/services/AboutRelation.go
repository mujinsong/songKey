package services

import (
	"log"
	"songKey/dao/graph"
	"songKey/domain"
)

func CreateRelation(relationCreate *domain.Relation) (*domain.Result, error) {
	cy := graph.CypherStruct{}
	cypher := cy.MatchNode(relationCreate.FromNode).MatchNode(relationCreate.ToNode).
		WhereAnd("n0", relationCreate.FromNode, []string{"id"}).WhereAnd("n1", relationCreate.ToNode, []string{"id"}).
		CreateOnlyRelation(relationCreate, "n0", "n1").
		ReturnAll().GetFinalCypher()
	log.Println("CreateRelation:" + cypher)
	res, err := cy.Result()
	return res, err
}
