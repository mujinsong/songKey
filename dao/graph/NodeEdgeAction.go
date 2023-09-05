package graph

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/global"
	"strings"
)

func AddNode(label string, properties []string, isUnique bool) (neo4j.Result, error) {
	var cypther strings.Builder
	cypther.WriteString("create (node:")
	cypther.WriteString(label)
	if len(properties) != 0 {
		cypther.WriteString("{")
		for _, property := range properties {
			cypther.WriteString(property)
			cypther.WriteString(":true,")
		}
		if isUnique {
			cypther.WriteString("unique:true")
		} else {
			cypther.WriteString("unique:false")
		}
		cypther.WriteString("}")
	}
	cypther.WriteString(")")
	result, err := Run(cypther.String())
	if err != nil {
		return result, err
	}
	return result, err
}

func Run(cypher string) (neo4j.Result, error) {
	session, err := global.Neo4jDriver.NewSession(neo4j.SessionConfig{})
	if err != nil {
		log.Fatalf("Failed to create Neo4j session: %v", err)
	}
	// 释放session
	defer session.Close()
	result, err := session.Run(cypher, nil)
	if err != nil {
		return nil, err
	}
	return result, err
}
func Exec(cypher string) (neo4j.Result, error) {
	session, err := global.Neo4jDriver.NewSession(neo4j.SessionConfig{})
	if err != nil {
		log.Fatalf("Failed to create Neo4j session: %v", err)
	}
	// 释放session
	defer session.Close()
	tx, err := session.BeginTransaction()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	result, err := tx.Run(cypher, nil)
	if err != nil {
		log.Fatalf("Failed to run query: %v", err)
	}
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
	return result, err
}
