package domain

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point. if you wanna new node, had better to use domain.NewNode
type Node struct {
	Id         int64             `json:"id" default:"-1"`
	IsUnique   bool              `json:"is_unique"`
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
	Relation
	IsDirect bool `json:"is_direct"`
	Min      int  `json:"min"`
	Max      int  `json:"max"`
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
