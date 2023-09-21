package domain

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"reflect"
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

	SetCypher strings.Builder

	WhereCypher     strings.Builder
	needConjunction bool
	ReturnCypher    strings.Builder
}

func (cypher *CypherStruct) Reset() {
	cypher.MatchCypher.Reset()
	cypher.matchCount = 0
	cypher.relationCount = 0
	cypher.needConjunction = false
	cypher.SetCypher.Reset()
	cypher.WhereCypher.Reset()
	cypher.ReturnCypher.Reset()
}
func (cypher *CypherStruct) SetNode(name string, node *Node) *CypherStruct {
	conjunction := ","
	if node.IsUnique {
		cypher.tryAddSetConjunction(conjunction)
		cypher.SetCypher.WriteString(fmt.Sprintf(" %s.unique=true ", name))
	}
	if len(node.Properties) > 0 {
		for pro, val := range node.Properties {
			cypher.tryAddSetConjunction(conjunction)
			cypher.SetCypher.WriteString(fmt.Sprintf(" %s.%s=%v ", name, pro, val))
		}
	}
	return cypher
}
func (cypher *CypherStruct) SetRelation(name string, relation *Relation) *CypherStruct {
	conjunction := ","
	if relation.ToNode.IsUnique || relation.ToNodeIsUnique {
		cypher.tryAddSetConjunction(conjunction)
		cypher.SetCypher.WriteString(fmt.Sprintf(" %s.unique=true ", name))
	}
	if len(relation.Properties) > 0 {
		for pro, val := range relation.Properties {
			cypher.tryAddSetConjunction(conjunction)
			cypher.SetCypher.WriteString(fmt.Sprintf(" %s.%s=%v ", name, pro, val))
		}
	}
	return cypher
}

// Set the `set` is a node or a relation, the `name` is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) Set(name string, set interface{}) *CypherStruct {
	theType := reflect.TypeOf(set)
	switch theType {
	case reflect.TypeOf(&Relation{}):
		cypher.SetRelation(name, (set).(*Relation))
		log.Println("SetRelation")
	case reflect.TypeOf(Relation{}):
		relation := (set).(Relation)
		cypher.SetRelation(name, &relation)
		log.Println("SetRelation")
	case reflect.TypeOf(&Node{}):
		cypher.SetNode(name, (set).(*Node))
		log.Println("SetNode")
	case reflect.TypeOf(Node{}):
		node := set.(Node)
		cypher.SetNode(name, &node)
		log.Println("SetNode")
	default:
		log.Println("Set: unKnow Type")
	}
	return cypher
}

// Where :if you use this function please write "where" by yourself
func (cypher *CypherStruct) Where(cy string) {
	cypher.WhereCypher.WriteString(cy)
}
func (cypher *CypherStruct) tryAddWhereConjunction(conjunction string) {
	if cypher.needConjunction {
		cypher.WhereCypher.WriteString(conjunction)
	} else {
		if cypher.WhereCypher.Len() == 0 {
			cypher.WhereCypher.WriteString(" where ")
		}
	}
}
func (cypher *CypherStruct) tryAddSetConjunction(conjunction string) {
	if cypher.needConjunction {
		cypher.SetCypher.WriteString(conjunction)
	} else {
		if cypher.SetCypher.Len() == 0 {
			cypher.SetCypher.WriteString(" set ")
		}
	}
}

// WhereRelation the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereRelation(name string, relation *Relation, conjunction string) *CypherStruct {
	if len(relation.Type) > 0 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" type(%s) =\"%s\"", name, relation.Type))
		cypher.needConjunction = true
	}
	if relation.ToNodeIsUnique {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.unique = \"true\"", name))
		cypher.needConjunction = true
	}
	if relation.Id != -1 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" id(%s) =%d", name, relation.Id))
		cypher.needConjunction = true
	}
	if len(relation.Properties) > 0 {
		for k, v := range relation.Properties {
			cypher.tryAddWhereConjunction(conjunction)
			cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.%s = \"%s\"", name, k, v))
			cypher.needConjunction = true
		}
	}
	return cypher
}

// WhereNode the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereNode(name string, node *Node, conjunction string) *CypherStruct {
	if len(node.Label) > 0 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" labels(%s) =[\"%s\"]", name, node.Label))
		cypher.needConjunction = true
	}
	if node.IsUnique {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.unique = \"true\"", name))
		cypher.needConjunction = true
	}
	if node.Id != -1 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" id(%s) =%d", name, node.Id))
		cypher.needConjunction = true
	}
	if len(node.Properties) > 0 {
		for k, v := range node.Properties {
			cypher.tryAddWhereConjunction(conjunction)
			cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.%s = \"%s\"", name, k, v))
			cypher.needConjunction = true
		}
	}
	return cypher
}

// WhereAnd :the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereAnd(name string, where interface{}) *CypherStruct {
	conjunction := " and "
	theType := reflect.TypeOf(where)
	switch theType {
	case reflect.TypeOf(&Relation{}):
		cypher.WhereRelation(name, (where).(*Relation), conjunction)
		log.Println("whereAndRelation")
	case reflect.TypeOf(Relation{}):
		relation := (where).(Relation)
		cypher.WhereRelation(name, &relation, conjunction)
		log.Println("whereAndRelation")
	case reflect.TypeOf(&Node{}):
		cypher.WhereNode(name, (where).(*Node), conjunction)
		log.Println("whereAndNode")
	case reflect.TypeOf(Node{}):
		node := where.(Node)
		cypher.WhereNode(name, &node, conjunction)
		log.Println("whereAndNode")

	default:
		log.Println("whereAnd: unKnow Type")
	}
	return cypher
}

// WhereOr :the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereOr(name string, where interface{}) *CypherStruct {
	conjunction := " or "
	theType := reflect.TypeOf(where)
	switch theType {
	case reflect.TypeOf(&Relation{}):
		cypher.WhereRelation(name, (where).(*Relation), conjunction)
		log.Println("whereOrRelation")
	case reflect.TypeOf(Relation{}):
		relation := (where).(Relation)
		cypher.WhereRelation(name, &relation, conjunction)
		log.Println("whereOrRelation")
	case reflect.TypeOf(&Node{}):
		cypher.WhereNode(name, (where).(*Node), conjunction)
		log.Println("whereOrNode")
	case reflect.TypeOf(Node{}):
		node := where.(Node)
		cypher.WhereNode(name, &node, conjunction)
		log.Println("whereOrNode")

	default:
		log.Println("whereOr: unKnow Type")
	}
	return cypher
}

func (cypher *CypherStruct) MatchNode(node *Node) *CypherStruct {
	isOk := false
	if node != nil && !utils.IsEmpty(node.Label) {
		isOk = cypher.MatchNodeByLabelStr(node.Label)
	}
	if !isOk {
		log.Println("No match Node!")
		if cypher.MatchCypher.Len() > 0 {
			cypher.MatchCypher.WriteByte(',')
		} else {
			cypher.MatchCypher.WriteString("match ")
		}
		cypher.MatchCypher.WriteString("()")
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

func (cypher *CypherStruct) concatRelationMatcher(relation *Relation) {
	cypher.MatchNode(relation.FromNode)
	cypher.MatchCypher.WriteByte('-')
	if relation.Type != "" {
		cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d:%s]", cypher.relationCount, relation.Type))
		cypher.relationCount++
	} else {
		cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d]", cypher.relationCount))
		cypher.relationCount++
	}
	cypher.MatchCypher.WriteString("->")
	if relation.ToNode != nil && !utils.IsEmpty(relation.ToNode.Label) {
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchCount, relation.ToNode.Label))
		cypher.matchCount++
	} else {
		cypher.MatchCypher.WriteString("()")
	}

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
	switch theType {
	case reflect.TypeOf(&Relation{}):
		log.Println("isMatcher")
		cypher.concatRelationMatcher(relation.(*Relation))
	case reflect.TypeOf(Relation{}):
		log.Println("isMatcher")
		temp := relation.(Relation)
		cypher.concatRelationMatcher(&temp)
	case reflect.TypeOf(&RelationQuery{}):
		log.Println("isQuery")
		cypher.concatRelationQuery(relation.(*RelationQuery))
	case reflect.TypeOf(RelationQuery{}):
		log.Println("isQuery")
		temp := relation.(RelationQuery)
		cypher.concatRelationQuery(&temp)
	default:
		log.Println("MatchRelation: unKnow Type")
	}
	return cypher
}

func (cypher *CypherStruct) ReturnNode() *CypherStruct {
	if cypher.matchCount <= 0 {
		return cypher
	}
	if cypher.ReturnCypher.Len() == 0 {
		cypher.ReturnCypher.WriteString(" return ")
	} else {
		cypher.ReturnCypher.WriteByte(',')
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
	if cypher.relationCount <= 0 {
		return cypher
	}
	if cypher.ReturnCypher.Len() == 0 {
		cypher.ReturnCypher.WriteString(" return ")
	} else {
		cypher.ReturnCypher.WriteByte(',')
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
	return cypher.MatchCypher.String() + " " + cypher.SetCypher.String() + " " + cypher.WhereCypher.String() + " " + cypher.ReturnCypher.String()
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
