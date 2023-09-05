package rds

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"songKey/global"
)

func ChangeDb(dbName string) (bool, error) {
	global.RdsDbName = dbName
	global.RDS_DSN = global.RdsUsername + ":" + global.RdsPassword + "@tcp(" + global.RdsUri + ")/" + global.RdsDbName + global.RdsOther
	err := ConnectRDS()
	if err != nil {
		log.Println("Change DB Fail")
		return false, err
	}
	return true, nil
}

func ConnectRDS() error {
	var err error = nil
	global.RdsDb, err = gorm.Open(mysql.Open(global.RDS_DSN), &gorm.Config{})
	if err != nil {
		log.Println("RDS Init Fail")
		return err
	}
	return err
}
