package graph

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/global"
)

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
func Exec(cypher string) {
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
	if result.Next() {
		count := result.Record().GetByIndex(0).(int64)
		fmt.Printf("Found %d nodes in the database", count)
	}
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}
