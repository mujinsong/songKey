package rds

import (
	"songKey/domain"
	"songKey/global"
)

func GetFieldMes(tableName string) []domain.FieldInfo {
	var fieldInfo []domain.FieldInfo
	global.RdsDb.Raw("desc " + tableName).Scan(&fieldInfo)
	return fieldInfo
}

func GetTableMes() []string {
	tableInfo := make([]string, 0)
	global.RdsDb.Raw("show tables").Scan(&tableInfo)
	return tableInfo
}
