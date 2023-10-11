package domain

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point. if you wanna new node, had better to use domain.NewNode
type Node struct {
	Id         int64             `json:"id" default:"-1"`
	IsUnique   bool              `json:"unique"`
	Label      string            `json:"label"`
	Properties map[string]string `json:"properties"`
}

func NewNode() *Node {
	return &Node{Id: -1, Properties: make(map[string]string)}
}
func (node *Node) Create() (*neo4j.Result, error) {
	return graph.CreateNode(node.Label, node.Properties, node.IsUnique)
}

// Relation :had better to use domain.NewRelation
type Relation struct {
	FromNode       *Node             `json:"from_node"`
	ToNode         *Node             `json:"to_node"`
	Id             int64             `json:"id" default:"-1"`
	ToNodeIsUnique bool              `json:"to_node_is_unique"`
	Type           string            `json:"type"`
	Properties     map[string]string `json:"properties"`
}

// RelationQuery :had better to use domain.NewRelationQuery
type RelationQuery struct {
	FromNode *Node `json:"from_node"`
	ToNode   *Node `json:"to_node"`
	IsDirect bool  `json:"is_direct" default:"true"`
	Min      int   `json:"min"`
	Max      int   `json:"max"`
}

func NewRelation() *Relation {
	return &Relation{Id: -1}
}
func NewRelationQuery() *RelationQuery {
	return &RelationQuery{FromNode: NewNode(), ToNode: NewNode(), IsDirect: true}
}

// Create relations
func (r *Relation) Create() (*neo4j.Result, error) {
	cypher := CypherStruct{}
	return cypher.CreateRelation(r).ReturnNode().ReturnRelation().Result()
}

type Result struct {
	Nodes     []Node
	Relations []Relation
}

func NeoResToResult(res *neo4j.Result) *Result {
	result := Result{Nodes: make([]Node, 0), Relations: make([]Relation, 0)}
	keys, err := (*res).Keys()
	if err != nil {
		log.Println("get keys error")
		return nil
	}
	for (*res).Next() {
		record := (*res).Record()
		for _, key := range keys {
			if value, ok := record.Get(key); ok {
				switch value.(type) {
				case neo4j.Node:
					node := value.(neo4j.Node)
					nd := Node{Id: node.Id(), Label: node.Labels()[0], Properties: make(map[string]string)}
					for k, v := range node.Props() {
						nd.Properties[k] = fmt.Sprintf("%v", v)
					}
					result.Nodes = append(result.Nodes, nd)
				case neo4j.Relationship:
					relation := value.(neo4j.Relationship)
					rls := Relation{Id: relation.Id(), Type: relation.Type(), Properties: make(map[string]string), ToNode: &Node{Id: relation.EndId()}, FromNode: &Node{Id: relation.StartId()}}
					for k, v := range relation.Props() {
						rls.Properties[k] = fmt.Sprintf("%v", v)
					}
					result.Relations = append(result.Relations, rls)
				}
			}
		}
	}
	return &result
}
