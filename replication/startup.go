package replication

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
)

var (
	hostname = flag.String("hostname", "", "Hostname of the node")
	main     = flag.Bool("main", false, "is this node the Main node")
)

type node struct {
	ID       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname"`
	Main     bool      `json:"main"`
}

type ReplicationController struct {
	EtcdClient *clientv3.Client
}

func (controller *ReplicationController) StartUp() error {
	if *hostname == "" {
		host, err := os.Hostname()
		if err != nil {
			return err
		}

		*hostname = host
	}

	n := &node{
		ID:       uuid.New(),
		Hostname: *hostname,
		Main:     *main,
	}

	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	_, err = controller.EtcdClient.Put(context.TODO(), fmt.Sprintf("node/%s", n.ID), string(data))
	if err != nil {
		return err
	}

	return nil
}
