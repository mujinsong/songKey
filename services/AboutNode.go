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
