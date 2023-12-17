package replication

import (
	"github.com/google/uuid"
)

type node struct {
	Id       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname"`
	Main     bool      `json:"main"`
}

func (n node) ID() string {
	return n.Id.String()
}

func (n node) Host() string {
	return n.Hostname
}

func (n node) IsMain() bool {
	return n.Main
}

type Node interface {
	ID() string
	Host() string
	IsMain() bool
}

var self *node

func GetSelf() Node {
	if self == nil {
		self = new(node)
		self.Id = uuid.New()
		self.Hostname = *hostname
		self.Main = *Main
	}

	return *self
}

type nodeService struct {
	nodes []*node
}

func (service *nodeService) GetMain() Node {
	for _, n := range service.nodes {
		if n.Main {
			return *n
		}
	}

	if len(service.nodes) == 0 {
		panic("no nodes registered")
	}

	return service.nodes[0]
}

func (service *nodeService) upsertNode(upsertNode *node) {
	for i, n := range service.nodes {
		if n.Id == upsertNode.Id {
			service.nodes[i] = upsertNode
			return
		}
	}

	service.nodes = append(service.nodes, upsertNode)
}
