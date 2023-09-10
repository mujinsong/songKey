package domain

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"reflect"
	"songKey/contants"
	"songKey/dao/graph"
	"songKey/utils"
	"strings"
	"sync"
)

type CypherStruct struct {
	MatchCypher   strings.Builder
	matchCount    int
	relationCount int
	matchLock     sync.Mutex

	ReturnCypher strings.Builder
	IsReturn     bool
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
		if cypher.matchCount != 0 {
			cypher.MatchCypher.WriteByte(',')
		} else {
			cypher.MatchCypher.WriteString("match ")
		}
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchCount, label))
		cypher.matchCount++
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
	cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d:%s]", cypher.matchCount, relation.Label))
	cypher.relationCount++
	cypher.MatchCypher.WriteString("->")
	cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchCount, relation.ToNode.Label))
	cypher.matchCount++

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
	cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchCount, relation.ToNode.Label))
	cypher.matchCount++
}

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

func (cypher *CypherStruct) ReturnNode() *CypherStruct {
	if !cypher.IsReturn {
		cypher.ReturnCypher.WriteString(" return ")
		cypher.IsReturn = true
	}

	for i := 0; i < cypher.matchCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("n%d", i))
		if i+1 < cypher.matchCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	return cypher
}
func (cypher *CypherStruct) ReturnRelation() *CypherStruct {
	if !cypher.IsReturn {
		cypher.ReturnCypher.WriteString(" return ")
		cypher.IsReturn = true
	}
	for i := 0; i < cypher.relationCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("r%d", i))
		if i+1 < cypher.relationCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	return cypher
}

func (cypher *CypherStruct) GetFinalCypher() string {
	return cypher.MatchCypher.String() + " " + cypher.ReturnCypher.String()
}

func (cypher *CypherStruct) Result() (*neo4j.Result, error) {
	finalCypher := cypher.GetFinalCypher()
	res, err := graph.Run(finalCypher)
	if err != nil {
		log.Println("error getResult")
		return nil, err
	}
	return res, err
}
