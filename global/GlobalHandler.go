package global

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"gorm.io/gorm"
)

var (
	InitResult  = false
	Neo4jDriver neo4j.Driver
	RdsDb       *gorm.DB
)
