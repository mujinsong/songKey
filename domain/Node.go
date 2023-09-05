package domain

type Node struct {
	NodeId   string `json:"id"`
	ObjId    string `json:"objId" `
	NodeName string `json:"node_name" `
	Desc     string `json:"desc"`
}
