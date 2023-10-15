package domain

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
	FromNode  *Node `json:"from_node"` //Deprecated: use StartNode instead if there are multiple nodes
	ToNode    *Node `json:"to_node"`   //Deprecated: use EndNode instead if there are multiple nodes
	IsDirect  bool  `json:"is_direct" default:"true"`
	Min       int   `json:"min"`
	Max       int   `json:"max"`
	StartNode []int `json:"start_node"`
	EndNode   []int `json:"end_node"`
}

func NewRelation() *Relation {
	return &Relation{Id: -1}
}
func NewRelationQuery() *RelationQuery {
	return &RelationQuery{FromNode: NewNode(), ToNode: NewNode(), IsDirect: true}
}

type Result struct {
	Nodes     []Node
	Relations []Relation
}
