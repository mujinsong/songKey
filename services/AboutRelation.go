package services

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/domain"
)

func CreateRelation(relationCreate *domain.Relation) (*neo4j.Result, error) {
	cy := domain.CypherStruct{}
	cypher := cy.MatchNode(relationCreate.FromNode).MatchNode(relationCreate.ToNode).WhereAnd("n0", relationCreate.FromNode).WhereAnd("n1", relationCreate.ToNode).CreateOnlyRelation(relationCreate, "n0", "n1").ReturnAll().GetFinalCypher()
	log.Println("CreateRelation:" + cypher)
	res, err := cy.Result()
	return res, err
}
