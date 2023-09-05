package main

import (
	"fmt"
	"songKey/dao/rds"
)

func test() {
	x := rds.GetTableMes()
	for _, str := range x {
		field := rds.GetFieldMes(str)
		fmt.Println(field)
	}
}
