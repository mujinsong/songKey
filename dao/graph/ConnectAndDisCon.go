package graph

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"songKey/global"
)

func ConnectNeo4j() error {
	var err error = nil
	global.Neo4jDriver, err = CreateDriver(global.Neo4jUri, global.Neo4jUsername, global.Neo4jPassword)
	if err != nil {
		log.Println("neo4j: init fail")
	}
	return err
}
func CreateDriver(uri, username, password string) (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
}
func CloseDriver(driver neo4j.Driver) error {
	return driver.Close()
}
