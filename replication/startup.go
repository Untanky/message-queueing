package replication

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log/slog"
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

	self     node
	nodeChan clientv3.WatchChan
	quit     chan bool
	err      chan error
}

func (controller *ReplicationController) StartUp(ctx context.Context) error {
	controller.quit = make(chan bool)
	controller.err = make(chan error)

	if *hostname == "" {
		host, err := os.Hostname()
		if err != nil {
			return err
		}

		*hostname = host
	}

	controller.self = node{
		ID:       uuid.New(),
		Hostname: *hostname,
		Main:     *main,
	}

	existingNodes, err := controller.getExistingNodes(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", existingNodes)

	slog.Info("registering this node", "nodeID", fmt.Sprintf("node/%s", controller.self.ID))
	err = controller.register(ctx)
	if err != nil {
		return err
	}

	controller.nodeChan = controller.EtcdClient.Watch(context.TODO(), fmt.Sprintf("node"), clientv3.WithPrefix())

	go controller.handleAsyncUpdates()

	return nil
}

func (controller *ReplicationController) register(ctx context.Context) error {
	data, err := json.Marshal(controller.self)
	if err != nil {
		return err
	}

	_, err = controller.EtcdClient.Put(ctx, fmt.Sprintf("node/%s", controller.self.ID), string(data))
	return err
}

func (controller *ReplicationController) getExistingNodes(ctx context.Context) ([]*node, error) {
	getResponse, err := controller.EtcdClient.Get(ctx, fmt.Sprintf("node"), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodes := make([]*node, 0, len(getResponse.Kvs))

	for _, kv := range getResponse.Kvs {
		n := new(node)
		err = json.Unmarshal(kv.Value, n)
		if err != nil {
			slog.Error("could not parse node information from etcd", "nodeID", kv.Key)
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (controller *ReplicationController) deregister(ctx context.Context) error {
	_, err := controller.EtcdClient.Delete(ctx, fmt.Sprintf("node/%s", controller.self.ID))
	return err
}

func (controller *ReplicationController) Close() error {
	controller.quit <- true
	return <-controller.err
}

func (controller *ReplicationController) handleAsyncUpdates() {
	for {
		select {
		case resp := <-controller.nodeChan:
			data := *resp.Events[0]
			if data.IsCreate() {
				slog.Info("got new node from etcd", "nodeID", data.Kv.Key)
			} else if data.IsModify() {
				slog.Info("updated node from etcd", "nodeID", data.Kv.Key)
			}
		case <-controller.quit:
			slog.Info("removing node data from etcd", "nodeID", controller.self.ID)

			err := controller.deregister(context.TODO())
			controller.err <- err
		}
	}
}
