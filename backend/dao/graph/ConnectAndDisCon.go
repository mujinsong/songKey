package graph

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	global2 "songKey/global"
)

func ConnectNeo4j() error {
	var err error = nil
	global2.Neo4jDriver, err = CreateDriver(global2.Neo4jUri, global2.Neo4jUsername, global2.Neo4jPassword)
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
