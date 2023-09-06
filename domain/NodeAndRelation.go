package domain

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"songKey/dao/graph"
)

// Node It is mandatory to have only one label for each point
type Node struct {
	isUnique   bool
	label      string
	properties map[string]string
}

func (node Node) Create() (neo4j.Result, error) {
	return graph.CreateNode(node.label, node.properties, node.isUnique)
}

type Relation struct {
	ToNodeIsUnique bool
	label          string
	properties     map[string]string
}
