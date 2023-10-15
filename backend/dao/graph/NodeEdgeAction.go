package graph

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"songKey/domain"
	"songKey/global"
	"strings"
)

func CreateNode(label string, properties map[string]string, isUnique bool) (*domain.Result, error) {
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
		return nil, err
	}
	return result, err
}

func Run(cypher string) (*domain.Result, error) {
	ctx := context.Background()
	session := global.Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{})
	result, err := session.Run(ctx, cypher, nil)
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Println("session-close-error:", err)
			recover()
		}
	}(session, ctx)
	if err != nil {
		return nil, err
	}
	return NeoResToResult(&result, true), err
}

//func Exec(cypher string) (neo4j.ResultWithContext, error) {
//	session, err := global.Neo4jDriver.NewSession(neo4j.SessionConfig{})
//	if err != nil {
//		log.Printf("Failed to create Neo4j session: %v\n", err)
//	}
//	// 释放session
//	defer session.Close()
//	tx, err := session.BeginTransaction()
//	if err != nil {
//		log.Printf("Failed to begin transaction: %v\n", err)
//	}
//	result, err := tx.Run(cypher, nil)
//	if err != nil {
//		log.Printf("Failed to run query: %v\n", err)
//	}
//	//提交事务
//	if err := tx.Commit; err != nil {
//		log.Printf("Failed to commit transaction: %v\n", err)
//	}
//	return result, err
//}

func NeoResToResult(res *neo4j.ResultWithContext, unique bool) *domain.Result {
	var uniqueNodeMp map[int64]bool = nil
	var uniquePathMp map[int64]bool = nil
	if unique {
		uniqueNodeMp = make(map[int64]bool)
		uniquePathMp = make(map[int64]bool)
	}
	result := domain.Result{Nodes: make([]domain.Node, 0), Relations: make([]domain.Relation, 0)}
	keys, err := (*res).Keys()
	if err != nil {
		log.Println("get keys error")
		return nil
	}
	ctx := context.Background()
	record := (*res).Record()
	for (*res).NextRecord(ctx, &record) {
		for _, key := range keys {
			if values, ok := record.Get(key); ok {
				for _, value := range values.([]interface{}) {
					switch value.(type) {
					case neo4j.Node:
						node := value.(neo4j.Node)
						addNode(uniqueNodeMp, &node, &result)
					case neo4j.Relationship:
						relation := value.(neo4j.Relationship)
						addRelationship(uniquePathMp, &relation, &result)
					case neo4j.Path:
						path := value.(neo4j.Path)
						for _, nd := range path.Nodes {
							addNode(uniqueNodeMp, &nd, &result)
						}
						for _, pt := range path.Relationships {
							addRelationship(uniqueNodeMp, &pt, &result)
						}
					}
				}
			}
		}
	}
	return &result
}

func addNode(uniqueNodeMp map[int64]bool, node *neo4j.Node, result *domain.Result) {
	if uniqueNodeMp != nil {
		if uniqueNodeMp[node.Id] {
			return
		} else {
			uniqueNodeMp[node.Id] = true
		}
	}
	nd := domain.Node{Id: node.Id, Label: node.Labels[0], Properties: make(map[string]string)}
	for k, v := range node.Props {
		nd.Properties[k] = fmt.Sprintf("%v", v)
	}
	result.Nodes = append(result.Nodes, nd)
}
func addRelationship(uniqueMp map[int64]bool, relation *neo4j.Relationship, result *domain.Result) {
	if uniqueMp != nil {
		if uniqueMp[relation.Id] {
			return
		} else {
			uniqueMp[relation.Id] = true
		}
	}
	rls := domain.Relation{Id: relation.Id, Type: relation.Type, Properties: make(map[string]string), ToNode: &domain.Node{Id: relation.EndId}, FromNode: &domain.Node{Id: relation.StartId}}
	for k, v := range relation.Props {
		rls.Properties[k] = fmt.Sprintf("%v", v)
	}
	result.Relations = append(result.Relations, rls)
}
