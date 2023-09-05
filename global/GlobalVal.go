package global

var (
	//for neo4j
	Neo4jUri      = "bolt://localhost:7687"
	Neo4jUsername = "graph"
	Neo4jPassword = "zsm20020609"

	//for RDS
	RdsUsername = "root"
	RdsPassword = "zsm20020609"
	RdsUri      = "127.0.0.1:3306"
	RdsDbName   = "singAndSong"
	RdsOther    = "?charset=utf8mb4&parseTime=True&loc=Local"
	RDS_DSN     = RdsUsername + ":" + RdsPassword + "@tcp(" + RdsUri + ")/" + RdsDbName + RdsOther
)
