package domain

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point
type Node struct {
	Id         int64
	IsUnique   bool
	Label      string
	Properties map[string]string
}

func (node *Node) Create() (*neo4j.Result, error) {
	return graph.CreateNode(node.Label, node.Properties, node.IsUnique)
}

type RelationMatcher struct {
	FromNode       *Node
	ToNode         *Node
	Id             int64
	ToNodeIsUnique bool
	Label          string
	Properties     map[string]string
}
type RelationQuery struct {
	FromNode       *Node
	ToNode         *Node
	Id             int64
	ToNodeIsUnique bool
	Label          string
	Properties     map[string]string
	IsDirect       bool
	Min            int
	Max            int
}
