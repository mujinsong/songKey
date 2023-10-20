package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"songKey/dao/graph"
	"songKey/dao/rds"
	"songKey/global"
	"songKey/manage"
)

func init() {
	err := rds.ConnectRDS()
	if err != nil {
		global.InitResult = false
		log.Println("connectRDS-error:", err)
		return
	}
	err = graph.ConnectNeo4j()
	if err != nil {
		global.InitResult = false
		log.Println("connectNeo4j-error", err)
		return
	}
	err = manage.InitGlobalMap()
	if err != nil {
		log.Println("initGlobalMap-error:", err)
		global.InitResult = false
		return
	}
	global.InitResult = true

}

func main() {
	if !global.InitResult {
		log.Fatalln("Init Fail")
		return
	}
	defer func(Neo4jDriver neo4j.DriverWithContext, ctx context.Context) {
		err := Neo4jDriver.Close(ctx)
		if err != nil {
			log.Println("Neo4jDriver close err:", err)
			recover()
		}
	}(global.Neo4jDriver, context.Background())
	StartRouter()
}
