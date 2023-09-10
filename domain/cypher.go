package domain

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"reflect"
	"songKey/contants"
	"songKey/utils"
	"strings"
	"sync"
)

type CypherStruct struct {
	MatchCypher   strings.Builder
	MatchCount    int
	RelationCount int
	matchLock     sync.Locker
}

func (cypher *CypherStruct) MatchNode(node *Node) *CypherStruct {
	isOk := false
	if node != nil && !utils.IsEmpty(node.Label) {
		isOk = cypher.MatchNodeByLabelStr(node.Label)
	}
	if !isOk {
		log.Println("Match Fail!")
	}
	return cypher
}
func (cypher *CypherStruct) MatchNodeByLabelStr(label string) bool {
	defer func() {
		if p := recover(); p != nil {
			cypher.matchLock.Unlock()
			log.Println("match Fail!!!")
		}
	}()
	if !utils.IsEmpty(label) {
		cypher.matchLock.Lock()
		if cypher.MatchCount != 0 {
			cypher.MatchCypher.WriteByte(',')
		} else {
			cypher.MatchCypher.WriteString("match ")
		}
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.MatchCount, label))
		cypher.MatchCount++
		cypher.matchLock.Unlock()
		return true
	} else {
		log.Println("Match Fail!")
		return false
	}
}

func (cypher *CypherStruct) concatRelationMatcher(relation *RelationMatcher) {
	cypher.MatchNode(relation.FromNode)
	cypher.MatchCypher.WriteByte('-')
	cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d:%s]", cypher.MatchCount, relation.Label))
	cypher.RelationCount++
	cypher.MatchCypher.WriteString("->")
	cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.MatchCount, relation.ToNode.Label))
	cypher.MatchCount++

}
func (cypher *CypherStruct) concatRelationQuery(relation *RelationQuery) {
	cypher.MatchNode(relation.FromNode)
	cypher.MatchCypher.WriteByte('-')
	if relation.IsDirect {
		cypher.MatchCypher.WriteString("[]->")
	} else {
		cypher.MatchCypher.WriteString("[*")
		if relation.Min > 0 && relation.Max <= 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("%d..", relation.Min))
		} else if relation.Min <= 0 && relation.Max > 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("..%d", relation.Max))
		} else if relation.Min > 0 && relation.Max > 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("%d..%d", relation.Min, relation.Max))
		}
		cypher.MatchCypher.WriteString("]->")
	}
	cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.MatchCount, relation.ToNode.Label))
	cypher.MatchCount++
}

// MatchRelation todo
func (cypher *CypherStruct) MatchRelation(relation interface{}) *CypherStruct {
	theType := reflect.TypeOf(relation)
	if theType.String() == contants.RELATION_MATCHER_NAME {
		log.Println("isMatcher")
		cypher.concatRelationMatcher(relation.(*RelationMatcher))
	} else if theType.String() == contants.RELATION_QUERY_NAME {
		log.Println("isQuery")
		cypher.concatRelationQuery(relation.(*RelationQuery))
	} else {
		log.Println("unKnow Type")
	}
	return cypher
}

// Return todo
func (cypher *CypherStruct) Return() *neo4j.Result {
	return nil
}
