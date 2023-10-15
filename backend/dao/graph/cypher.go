package graph

import (
	"fmt"
	"log"
	"songKey/domain"
	"songKey/myUtils"
	"strings"
	"sync"
)

type CypherStruct struct {
	MatchCypher        strings.Builder
	matchNodeCount     int
	matchRelationCount int
	matchLock          sync.Mutex

	createCypher        strings.Builder
	createNodeCount     int
	createRelationCount int

	SetCypher strings.Builder

	WhereCypher     strings.Builder
	needConjunction bool
	ReturnCypher    strings.Builder
}

func (cypher *CypherStruct) Reset() {
	cypher.MatchCypher.Reset()
	cypher.matchNodeCount = 0
	cypher.matchRelationCount = 0
	cypher.needConjunction = false
	cypher.createNodeCount = 0
	cypher.createRelationCount = 0
	cypher.createCypher.Reset()
	cypher.SetCypher.Reset()
	cypher.WhereCypher.Reset()
	cypher.ReturnCypher.Reset()
}

func (cypher *CypherStruct) SetNode(name string, node *domain.Node, mp map[string]bool) *CypherStruct {
	conjunction := ","
	if mp["unique"] || mp["all"] {
		if node.IsUnique {
			cypher.tryAddSetConjunction(conjunction)
			cypher.SetCypher.WriteString(fmt.Sprintf(" %s.unique=true ", name))
		} else {
			cypher.tryAddSetConjunction(conjunction)
			cypher.SetCypher.WriteString(fmt.Sprintf(" %s.unique=false ", name))
		}
	}
	if len(node.Properties) > 0 {
		for pro, val := range node.Properties {
			if mp[pro] || mp["all"] {
				cypher.tryAddSetConjunction(conjunction)
				cypher.SetCypher.WriteString(fmt.Sprintf(" %s.%s=%v ", name, pro, val))
			}
		}
	}
	return cypher
}
func (cypher *CypherStruct) SetRelation(name string, relation *domain.Relation) *CypherStruct {
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
func (cypher *CypherStruct) Set(name string, set interface{}, element []string) *CypherStruct {
	mp := make(map[string]bool)
	for _, v := range element {
		mp[v] = true
	}
	switch set.(type) {
	case *domain.Relation:
		cypher.SetRelation(name, (set).(*domain.Relation))
		log.Println("SetRelation")
	case domain.Relation:
		relation := (set).(domain.Relation)
		cypher.SetRelation(name, &relation)
		log.Println("SetRelation")
	case *domain.Node:
		cypher.SetNode(name, (set).(*domain.Node), mp)
		log.Println("SetNode")
	case domain.Node:
		node := set.(domain.Node)
		cypher.SetNode(name, &node, mp)
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
func (cypher *CypherStruct) WhereRelation(name string, relation *domain.Relation, conjunction string, mp map[string]bool) *CypherStruct {
	if len(relation.Type) > 0 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" type(%s) =\"%s\"", name, relation.Type))
		cypher.needConjunction = true
	}
	if relation.ToNodeIsUnique && (mp["unique"] || mp["all"]) {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.unique = \"true\"", name))
		cypher.needConjunction = true
	}
	if relation.Id != -1 && (mp["id"] || mp["all"]) {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" id(%s) =%d", name, relation.Id))
		cypher.needConjunction = true
	}
	if len(relation.Properties) > 0 {
		for k, v := range relation.Properties {
			if mp[k] || mp["all"] {
				cypher.tryAddWhereConjunction(conjunction)
				cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.%s = \"%s\"", name, k, v))
				cypher.needConjunction = true
			}
		}
	}
	return cypher
}

// WhereNode the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereNode(name string, node *domain.Node, conjunction string, element map[string]bool) *CypherStruct {
	if len(node.Label) > 0 {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" labels(%s) =[\"%s\"]", name, node.Label))
		cypher.needConjunction = true
	}
	if node.IsUnique && (element["unique"] || element["all"]) {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.unique = \"true\"", name))
		cypher.needConjunction = true
	}
	if node.Id != -1 && (element["id"] || element["all"]) {
		cypher.tryAddWhereConjunction(conjunction)
		cypher.WhereCypher.WriteString(fmt.Sprintf(" id(%s) =%d", name, node.Id))
		cypher.needConjunction = true
	}
	if len(node.Properties) > 0 {
		for k, v := range node.Properties {
			if element[k] || element["all"] {
				cypher.tryAddWhereConjunction(conjunction)
				cypher.WhereCypher.WriteString(fmt.Sprintf(" %s.%s = \"%s\"", name, k, v))
				cypher.needConjunction = true
			}
		}
	}
	return cypher
}

// WhereAnd :the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereAnd(name string, where interface{}, element []string) *CypherStruct {
	conjunction := " and "
	mp := make(map[string]bool)
	for _, v := range element {
		mp[v] = true
	}
	switch where.(type) {
	case *domain.Relation:
		cypher.WhereRelation(name, (where).(*domain.Relation), conjunction, mp)
		log.Println("whereAndRelation")
	case domain.Relation:
		relation := (where).(domain.Relation)
		cypher.WhereRelation(name, &relation, conjunction, mp)
		log.Println("whereAndRelation")
	case *domain.Node:
		cypher.WhereNode(name, (where).(*domain.Node), conjunction, mp)
		log.Println("whereAndNode")
	case domain.Node:
		node := where.(domain.Node)
		cypher.WhereNode(name, &node, conjunction, mp)
		log.Println("whereAndNode")

	default:
		log.Println("whereAnd: unKnow Type")
	}
	return cypher
}

// WhereOr :the name is the param of match, like n1,n2(is node),r1,r2(is relationship)
func (cypher *CypherStruct) WhereOr(name string, where interface{}, element []string) *CypherStruct {
	conjunction := " or "
	mp := make(map[string]bool)
	for _, v := range element {
		mp[v] = true
	}
	switch where.(type) {
	case *domain.Relation:
		cypher.WhereRelation(name, (where).(*domain.Relation), conjunction, mp)
		log.Println("whereOrRelation")
	case domain.Relation:
		relation := (where).(domain.Relation)
		cypher.WhereRelation(name, &relation, conjunction, mp)
		log.Println("whereOrRelation")
	case *domain.Node:
		cypher.WhereNode(name, (where).(*domain.Node), conjunction, mp)
		log.Println("whereOrNode")
	case domain.Node:
		node := where.(domain.Node)
		cypher.WhereNode(name, &node, conjunction, mp)
		log.Println("whereOrNode")

	default:
		log.Println("whereOr: unKnow Type")
	}
	return cypher
}

func (cypher *CypherStruct) MatchNode(node *domain.Node) *CypherStruct {
	isOk := false
	if node != nil && !myUtils.IsEmpty(node.Label) {
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
func (cypher *CypherStruct) CreateNode(node *domain.Node) *CypherStruct {
	isOk := false
	if node != nil && !myUtils.IsEmpty(node.Label) {
		isOk = cypher.CreateNodeByLabelStr(node.Label, node.Properties, node.IsUnique)
	}
	if !isOk {
		log.Println("No create Node!")
		if cypher.createCypher.Len() > 0 {
			cypher.createCypher.WriteByte(',')
		} else {
			cypher.createCypher.WriteString("create ")
		}
		cypher.createCypher.WriteString("()")
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
	if !myUtils.IsEmpty(label) {
		cypher.matchLock.Lock()
		if cypher.matchNodeCount != 0 {
			cypher.MatchCypher.WriteByte(',')
		} else {
			cypher.MatchCypher.WriteString("match ")
		}
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchNodeCount, label))
		cypher.matchNodeCount++
		cypher.matchLock.Unlock()
		return true
	} else {
		log.Println("Match Fail!")
		return false
	}
}
func (cypher *CypherStruct) CreateNodeByLabelStr(label string, properties map[string]string, isUnique bool) bool {
	defer func() {
		if p := recover(); p != nil {
			cypher.matchLock.Unlock()
			log.Println("match Fail!!!")
		}
	}()
	if !myUtils.IsEmpty(label) {
		cypher.matchLock.Lock()
		if cypher.createNodeCount != 0 {
			cypher.createCypher.WriteByte(',')
		} else {
			cypher.createCypher.WriteString("create ")
		}
		cypher.createCypher.WriteString(fmt.Sprintf("(n%d:%s", cypher.createNodeCount, label))
		cypher.createNodeCount++
		cypher.createCypher.WriteString("{")
		if len(properties) > 0 {
			for k, v := range properties {
				cypher.createCypher.WriteString(fmt.Sprintf("%s:%s,", k, v))
			}
		}
		if isUnique {
			cypher.createCypher.WriteString(fmt.Sprintf("unique:true})"))
		} else {
			cypher.createCypher.WriteString(fmt.Sprintf("unique:false})"))
		}
		cypher.matchLock.Unlock()
		return true
	} else {
		log.Println("create Fail!")
		return false
	}
}

func (cypher *CypherStruct) concatRelationMatcher(relation *domain.Relation) *CypherStruct {
	cypher.MatchNode(relation.FromNode)
	cypher.MatchCypher.WriteByte('-')
	if relation.Type != "" {
		cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d:%s]", cypher.matchRelationCount, relation.Type))
		cypher.matchRelationCount++
	} else {
		cypher.MatchCypher.WriteString(fmt.Sprintf("[r%d]", cypher.matchRelationCount))
		cypher.matchRelationCount++
	}
	cypher.MatchCypher.WriteString("->")
	if relation.ToNode != nil && !myUtils.IsEmpty(relation.ToNode.Label) {
		cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchNodeCount, relation.ToNode.Label))
		cypher.matchNodeCount++
	} else {
		cypher.MatchCypher.WriteString("()")
	}
	return cypher
}
func (cypher *CypherStruct) CreateOnlyRelation(relation *domain.Relation, fromNode, toNode string) *CypherStruct {
	if cypher.createCypher.Len() == 0 {
		cypher.createCypher.WriteString(" create ")
	} else {
		cypher.createCypher.WriteByte(',')
	}
	cypher.createCypher.WriteString(fmt.Sprintf("(%s)-[r%d:%s{", fromNode, cypher.createRelationCount, relation.Type))
	cypher.createRelationCount++
	for k, pro := range relation.Properties {
		cypher.createCypher.WriteString(fmt.Sprintf("%s:%v,", k, pro))
	}
	if relation.ToNodeIsUnique == true {
		cypher.createCypher.WriteString("unique:true}")
	} else {
		cypher.createCypher.WriteString("unique:false}")
	}
	cypher.createCypher.WriteString(fmt.Sprintf("]->(%s) ", toNode))
	return cypher
}
func (cypher *CypherStruct) concatRelationCreate(relation *domain.Relation) *CypherStruct {
	cypher.CreateNode(relation.FromNode)
	cypher.createCypher.WriteByte('-')
	if relation.Type != "" {
		cypher.createCypher.WriteString(fmt.Sprintf("[r%d:%s]", cypher.createRelationCount, relation.Type))
		cypher.createRelationCount++
	} else {
		cypher.createCypher.WriteString(fmt.Sprintf("[r%d]", cypher.createRelationCount))
		cypher.createRelationCount++
	}
	cypher.createCypher.WriteString("->")
	if relation.ToNode != nil && !myUtils.IsEmpty(relation.ToNode.Label) {
		cypher.createCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.createNodeCount, relation.ToNode.Label))
		cypher.createNodeCount++
	} else {
		cypher.createCypher.WriteString("()")
	}
	return cypher
}
func (cypher *CypherStruct) concatRelationQuery(relation *domain.RelationQuery) *CypherStruct {
	cypher.MatchNode(relation.FromNode)
	cypher.MatchCypher.WriteByte('-')
	if relation.IsDirect {
		cypher.MatchCypher.WriteString("[r0]->")
	} else {
		cypher.MatchCypher.WriteString("[r0*")
		if relation.Min > 0 && relation.Max <= 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("%d..", relation.Min))
		} else if relation.Min <= 0 && relation.Max > 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("..%d", relation.Max))
		} else if relation.Min > 0 && relation.Max > 0 {
			cypher.MatchCypher.WriteString(fmt.Sprintf("%d..%d", relation.Min, relation.Max))
		}
		cypher.MatchCypher.WriteString("]->")
	}
	cypher.matchRelationCount++
	cypher.MatchCypher.WriteString(fmt.Sprintf("(n%d:%s)", cypher.matchNodeCount, relation.ToNode.Label))
	cypher.matchNodeCount++
	return cypher
}

func (cypher *CypherStruct) MatchRelation(relation interface{}) *CypherStruct {
	switch relation.(type) {
	case *domain.Relation:
		log.Println("isMatcher")
		cypher.concatRelationMatcher(relation.(*domain.Relation))
	case domain.Relation:
		log.Println("isMatcher")
		temp := relation.(domain.Relation)
		cypher.concatRelationMatcher(&temp)
	case *domain.RelationQuery:
		log.Println("isQuery")
		cypher.concatRelationQuery(relation.(*domain.RelationQuery))
	case domain.RelationQuery:
		log.Println("isQuery")
		temp := relation.(domain.RelationQuery)
		cypher.concatRelationQuery(&temp)
	default:
		log.Println("MatchRelation: unKnow Type")
	}
	return cypher
}
func (cypher *CypherStruct) CreateRelation(relation *domain.Relation) *CypherStruct {
	cypher.concatRelationCreate(relation)
	return cypher
}

func (cypher *CypherStruct) ReturnNode() *CypherStruct {
	if cypher.matchNodeCount <= 0 && cypher.createNodeCount <= 0 {
		return cypher
	}
	if cypher.ReturnCypher.Len() == 0 {
		cypher.ReturnCypher.WriteString(" return ")
	} else {
		cypher.ReturnCypher.WriteByte(',')
	}

	for i := 0; i < cypher.matchNodeCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("n%d", i))
		if i+1 < cypher.matchNodeCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	for i := 0; i < cypher.createNodeCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("n%d", i))
		if i+1 < cypher.createNodeCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	return cypher
}
func (cypher *CypherStruct) ReturnRelation() *CypherStruct {
	if cypher.matchRelationCount <= 0 && cypher.createRelationCount <= 0 {
		return cypher
	}
	if cypher.ReturnCypher.Len() == 0 {
		cypher.ReturnCypher.WriteString(" return ")
	} else {
		cypher.ReturnCypher.WriteByte(',')
	}
	for i := 0; i < cypher.matchRelationCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("r%d", i))
		if i+1 < cypher.matchRelationCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	for i := 0; i < cypher.createRelationCount; i++ {
		cypher.ReturnCypher.WriteString(fmt.Sprintf("r%d", i))
		if i+1 < cypher.createRelationCount {
			cypher.ReturnCypher.WriteByte(',')
		}
	}
	return cypher
}
func (cypher *CypherStruct) ReturnAll() *CypherStruct {
	cypher.ReturnCypher.WriteString(" return *")
	return cypher
}

func (cypher *CypherStruct) GetFinalCypher() string {
	return cypher.MatchCypher.String() + " " + cypher.WhereCypher.String() + " " + cypher.SetCypher.String() + " " + cypher.createCypher.String() + " " + cypher.ReturnCypher.String()
}

func (cypher *CypherStruct) Result() (*domain.Result, error) {
	finalCypher := cypher.GetFinalCypher()
	res, err := Run(finalCypher)
	if err != nil {
		log.Println("error getResult")
		return nil, err
	}
	return res, err
}
