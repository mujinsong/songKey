package domain

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point
type Node struct {
	Id         int64 `default:"-1"`
	IsUnique   bool
	Label      string
	Properties map[string]string
}

func NewNode() *Node {
	return &Node{Id: -1, Properties: make(map[string]string)}
}
func (node *Node) Create() (*neo4j.Result, error) {
	return graph.CreateNode(node.Label, node.Properties, node.IsUnique)
}

type Relation struct {
	FromNode       *Node
	ToNode         *Node
	Id             int64 `default:"-1"`
	ToNodeIsUnique bool
	Type           string
	Properties     map[string]string
}
type RelationQuery struct {
	Relation
	IsDirect bool
	Min      int
	Max      int
}

func NewRelation() *Relation {
	return &Relation{Id: -1}
}
func NewRelationQuery() *RelationQuery {
	return &RelationQuery{Relation: *NewRelation()}
}

// Create relations
func (r *Relation) Create() (*neo4j.Result, error) {
	cypher := CypherStruct{}
	return cypher.CreateRelation(r).ReturnNode().ReturnRelation().Result()
}
