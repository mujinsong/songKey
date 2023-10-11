package services

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/domain"
)

// todo
func QueryPath(query *domain.RelationQuery) (*neo4j.Result, error) {
	cy := domain.CypherStruct{}
	cypher := cy.MatchRelation(query).ReturnAll().GetFinalCypher()
	log.Println("service-QueryPath-cypher:", cypher)
	return nil, nil
}
