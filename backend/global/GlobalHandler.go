package global

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

var (
	InitResult                          = false
	Neo4jDriver neo4j.DriverWithContext = nil
	RdsDb       *gorm.DB                = nil
	KVMap       map[string]interface{}  = nil
)
