package manage

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"log"
	"os"
	"songKey/global"
)

func InitGlobalMap() error {
	global.KVMap = make(map[string]interface{})
	byteValue, err := os.ReadFile("backend/resource/CypherMap.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(byteValue, &global.KVMap)
	if err != nil {
		log.Println("initGlobalMap-error:", err)
		return err
	}
	return nil
}
