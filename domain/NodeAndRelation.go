package domain

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point
type Node struct {
	Id         int64 `default:"-1"`
	IsUnique   bool  `default:"false"`
	Label      string
	Properties map[string]string
}

func (node *Node) Create() (*neo4j.Result, error) {
	return graph.CreateNode(node.Label, node.Properties, node.IsUnique)
}

// Modify todo
func (node *Node) Modify() (*neo4j.Result, error) {
	return nil, nil
}

type Relation struct {
	FromNode       *Node
	ToNode         *Node
	Id             int64 `default:"-1"`
	ToNodeIsUnique bool
	Label          string
	Properties     map[string]string
}
type RelationQuery struct {
	Relation
	IsDirect bool
	Min      int
	Max      int
}

// Create todo
func (r *Relation) Create() (*neo4j.Result, error) {

	return nil, nil
}
