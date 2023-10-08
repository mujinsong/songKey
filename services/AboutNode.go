package services

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/domain"
)

func CreateNode(node *domain.Node) (*neo4j.Result, error) {
	cy := domain.CypherStruct{}
	cypher := cy.CreateNode(node).ReturnNode().GetFinalCypher()
	log.Println("service-createNode:" + cypher)
	res, err := cy.Result()
	return res, err
}

func SetNode(node *domain.Node) (*neo4j.Result, error) {
	cy := domain.CypherStruct{}
	cypher := cy.MatchNode(node).Set("n0", node, []string{"all"}).WhereAnd("n0", node, []string{"id"}).ReturnNode().GetFinalCypher()
	log.Println("service-SetNode: " + cypher)
	result, err := cy.Result()
	if err != nil {
		log.Println("service-setNode-err:", err)
		return nil, err
	}
	return result, err
}
