package main

import (
	"log"
	"songKey/dao/graph"
	"songKey/dao/rds"
	"songKey/global"
)

func init() {
	err := rds.ConnectRDS()
	if err != nil {
		global.InitResult = false
		return
	}
	err = graph.ConnectNeo4j()
	if err != nil {
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
	test()
}
