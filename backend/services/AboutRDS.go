package services

import (
	"songKey/dao/rds"
	"songKey/domain"
)

func GetFieldMes(tables []string) map[string][]domain.FieldInfo {
	result := make(map[string][]domain.FieldInfo)
	for _, table := range tables {
		result[table] = rds.GetFieldMes(table)
	}
	return result
}
