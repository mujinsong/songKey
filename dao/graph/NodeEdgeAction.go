package graph

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/global"
	"strings"
)

func CreateNode(label string, properties map[string]string, isUnique bool) (*neo4j.Result, error) {
	var cypher strings.Builder
	cypher.WriteString("create (node:")
	cypher.WriteString(label)
	if len(properties) != 0 {
		cypher.WriteString("{")
		for k, property := range properties {
			cypher.WriteString(k)
			cypher.WriteString(":" + property + ",")
		}
		if isUnique {
			cypher.WriteString("unique:true")
		} else {
			cypher.WriteString("unique:false")
		}
		cypher.WriteString("}")
	}
	cypher.WriteString(")")
	result, err := Run(cypher.String())
	if err != nil {
		return result, err
	}
	return result, err
}

func Run(cypher string) (*neo4j.Result, error) {
	session, err := global.Neo4jDriver.NewSession(neo4j.SessionConfig{})
	if err != nil {
		log.Printf("Failed to create Neo4j session: %v\n", err)
	}
	// 释放session
	defer session.Close()
	result, err := session.Run(cypher, nil)
	if err != nil {
		return nil, err
	}
	return &result, err
}
func Exec(cypher string) (neo4j.Result, error) {
	session, err := global.Neo4jDriver.NewSession(neo4j.SessionConfig{})
	if err != nil {
		log.Printf("Failed to create Neo4j session: %v\n", err)
	}
	// 释放session
	defer session.Close()
	tx, err := session.BeginTransaction()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
	}
	result, err := tx.Run(cypher, nil)
	if err != nil {
		log.Printf("Failed to run query: %v\n", err)
	}
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
	}
	return result, err
}
