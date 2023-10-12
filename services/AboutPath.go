package services

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/domain"
)

// QueryPath todo: handle the problem from issue
func QueryPath(query *domain.RelationQuery) (*neo4j.Result, error) {
	cy := domain.CypherStruct{}
	cy.MatchRelation(query).WhereAnd("n0", query.FromNode, []string{"id"}).WhereAnd("n1", query.ToNode, []string{"id"})
	cypher := cy.ReturnAll().GetFinalCypher()
	log.Println("service-QueryPath-cypher:", cypher)
	result, err := cy.Result()
	if err != nil {
		log.Println("service-queryPath-error:", err)
		return nil, err
	}
	return result, nil
}
