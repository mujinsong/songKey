package services

import (
	"fmt"
	"log"
	"songKey/dao/graph"
	"songKey/domain"
)

func CreateNode(node *domain.Node) (*domain.Result, error) {
	cy := graph.CypherStruct{}
	cypher := cy.CreateNode(node).ReturnNode().GetFinalCypher()
	log.Println("service-createNode:" + cypher)
	res, err := cy.Result()
	return res, err
}

func SetNode(node *domain.Node) (*domain.Result, error) {
	cy := graph.CypherStruct{}
	cypher := cy.MatchNode(node).Set("n0", node, []string{"all"}).WhereAnd("n0", node, []string{"id"}).ReturnNode().GetFinalCypher()
	log.Println("service-SetNode: " + cypher)
	result, err := cy.Result()
	if err != nil {
		log.Println("service-setNode-err:", err)
		return nil, err
	}
	return result, err
}
func MatchNode(node *domain.Node) (*domain.Result, error) {
	cy := graph.CypherStruct{}
	cypher := cy.MatchNode(node).WhereAnd("n0", node, []string{"all"}).ReturnNode().GetFinalCypher()
	log.Println("service-MatchNode: " + cypher)
	result, err := cy.Result()
	if err != nil {
		log.Println("service-MatchNode-err:", err)
		return nil, err
	}
	return result, err
}

func DeleteNode(nodes []*domain.Node) (*domain.Result, error) {
	cy := graph.CypherStruct{}
	for i, node := range nodes {
		cy.MatchNode(node).WhereAnd(fmt.Sprintf("n%d", i), node, []string{"all"}).Delete(fmt.Sprintf("n%d", i))
	}
	cypher := cy.ReturnAll().GetFinalCypher()
	log.Println("service-MatchNode: " + cypher)
	result, err := cy.Result()
	if err != nil {
		log.Println("service-DeleteNodes-err:", err)
		return nil, err
	}
	return result, err
}
