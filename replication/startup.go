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
}

func (controller *ReplicationController) StartUp() error {
	controller.quit = make(chan bool)

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

	data, err := json.Marshal(controller.self)
	if err != nil {
		return err
	}

	_, err = controller.EtcdClient.Get(context.TODO(), fmt.Sprintf("node"), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	controller.nodeChan = controller.EtcdClient.Watch(context.TODO(), fmt.Sprintf("node"), clientv3.WithPrefix())

	_, err = controller.EtcdClient.Put(context.TODO(), fmt.Sprintf("node/%s", controller.self.ID), string(data))
	if err != nil {
		return err
	}

	go controller.handleAsyncUpdates()

	return nil
}

func (controller *ReplicationController) Close() error {
	controller.quit <- true
	return nil
}

func (controller *ReplicationController) handleAsyncUpdates() {
	for {
		select {
		case resp := <-controller.nodeChan:
			data := *resp.Events[0]
			slog.Info("got new node from etcd", "nodeID", data.Kv.Key)
		case <-controller.quit:
			slog.Info("removing node data from etcd", "nodeID", controller.self.ID)

			_, err := controller.EtcdClient.Delete(context.TODO(), fmt.Sprintf("node/%s", controller.self.ID))
			if err != nil {
				panic(err)
			}
		}
	}
}
