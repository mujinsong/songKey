package services

import (
	"fmt"
	"log"
	"os"
	"songKey/dao/graph"
	"songKey/domain"
	"songKey/global"
	"strconv"
	"strings"
)

// QueryPath :query must have startNode and EndNode,FromNode and ToNode are irrelevant
func QueryPath(query *domain.RelationQuery) (*domain.Result, error) {
	startNode := strings.Builder{}
	endNode := strings.Builder{}
	startNode.WriteByte('[')
	for i := 0; i < len(query.StartNode)-1; i++ {
		startNode.WriteString(fmt.Sprintf("%d,", query.StartNode[i]))
	}
	startNode.WriteString(fmt.Sprintf("%d]", query.StartNode[len(query.StartNode)-1]))
	endNode.WriteByte('[')
	for i := 0; i < len(query.EndNode)-1; i++ {
		endNode.WriteString(fmt.Sprintf("%d,", query.EndNode[i]))
	}
	endNode.WriteString(fmt.Sprintf("%d]", query.EndNode[len(query.EndNode)-1]))
	scope := ""

	if query.Min != 0 {
		scope = fmt.Sprintf("%d..", query.Min)
	}
	if query.Max != 0 && query.Min <= query.Max {
		if len(scope) > 0 {
			scope += strconv.Itoa(query.Max)
		} else {
			scope = fmt.Sprintf("..%d", query.Max)
		}
	}

	cypher := fmt.Sprintf(getCypherAboutQueryPath(), startNode.String(), endNode.String(), scope)
	log.Println(cypher)
	res, err := graph.Run(cypher)
	if err != nil {
		log.Println("service-queryPath-error:", err)
		return nil, err
	}
	return res, nil
}

func getCypherAboutQueryPath() string {
	cypherMap := global.KVMap["cypher"].(map[string]interface{})
	f, err := os.ReadFile(cypherMap["QueryPath"].(string))
	if err != nil {
		log.Println("service-getCypherAboutQueryPath-get-cypher-String-error:", err)
		return ""
	}
	return string(f)
}
