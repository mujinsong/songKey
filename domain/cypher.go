package domain

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/utils"
	"strings"
)

type CypherStruct struct {
	MatchCypher strings.Builder
	MatchCount  int
}

func (cypher *CypherStruct) MatchNode(node *Node) bool {
	if node != nil && !utils.IsEmpty(node.label) {
		return cypher.MatchNodeByLabelStr(node.label)
	} else {
		log.Println("Match Fail!")
		return false
	}
}
func (cypher *CypherStruct) MatchNodeByLabelStr(label string) bool {
	if !utils.IsEmpty(label) {
		if cypher.MatchCount != 0 {
			cypher.MatchCypher.WriteByte(',')
		} else {
			cypher.MatchCypher.WriteString("match ")
		}
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.MatchCount, label))
		cypher.MatchCount++
		return true
	} else {
		log.Println("Match Fail!")
		return false
	}
}

// MatchRelation todo
func (cypher *CypherStruct) MatchRelation(fromNode *Node) {
	cypher.MatchNode(fromNode)
}

// Return todo
func (cypher *CypherStruct) Return() *neo4j.Result {
	return nil
}
