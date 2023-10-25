package myUtils

import (
	"github.com/cloudwego/hertz/pkg/common/json"
	"log"
	"songKey/domain"
)

func IsEmpty(str string) bool {
	if &str == nil {
		return true
	}
	if len(str) == 0 {
		return true
	}
	if str == "" {
		return true
	}
	return false
}
func NodesGet(body []byte) ([]*domain.Node, error) {
	temp := make([]domain.Node, 1)
	err := json.Unmarshal(body, &temp)
	if err != nil {
		log.Println("SetNode-unmarshal-err:", err)
		return nil, err
	}
	l := len(temp)
	nodes := make([]*domain.Node, l)
	for i := 0; i < l; i++ {
		nodes[i] = domain.NewNode()
	}
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Println("SetNode-unmarshal-err:", err)
		return nil, err
	}
	return nodes, nil
}
