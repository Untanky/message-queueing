package replication

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"log/slog"
	http2 "message-queueing/http"
	"net/http"
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

type Controller struct {
	etcdClient *clientv3.Client
	self       node
	nodeChan   clientv3.WatchChan
	quit       chan bool
	err        chan error
	closed     bool
}

func Open(ctx context.Context, etcdClient *clientv3.Client) (*Controller, error) {
	quit := make(chan bool)
	errChan := make(chan error)

	if *hostname == "" {
		host, err := os.Hostname()
		if err != nil {
			return nil, err
		}

		*hostname = host
	}

	self := node{
		ID:       uuid.New(),
		Hostname: *hostname,
		Main:     *main,
	}

	controller := Controller{
		quit:       quit,
		err:        errChan,
		self:       self,
		etcdClient: etcdClient,
		nodeChan:   etcdClient.Watch(context.TODO(), fmt.Sprintf("node"), clientv3.WithPrefix()),
		closed:     false,
	}

	if !controller.self.Main {
		existingNodes, err := controller.getExistingNodes(ctx)
		if err != nil {
			return nil, err
		}

		var mainNode *node
		for _, existingNode := range existingNodes {
			if existingNode.Main {
				mainNode = existingNode
			}
		}

		if mainNode == nil {
			return nil, err
		}

		err = controller.syncFiles(ctx, mainNode)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("registering this node", "nodeID", fmt.Sprintf("node/%s", controller.self.ID))
	err := controller.register(ctx)
	if err != nil {
		return nil, err
	}

	go controller.handleAsyncUpdates()

	return &Controller{}, nil
}

func (controller *Controller) register(ctx context.Context) error {
	data, err := json.Marshal(controller.self)
	if err != nil {
		return err
	}

	_, err = controller.etcdClient.Put(ctx, fmt.Sprintf("node/%s", controller.self.ID), string(data))
	return err
}

func (controller *Controller) getExistingNodes(ctx context.Context) ([]*node, error) {
	getResponse, err := controller.etcdClient.Get(ctx, fmt.Sprintf("node"), clientv3.WithPrefix())
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

func (mainNode *node) getManifest(ctx context.Context) (*http2.GetManifestResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/internal/queue/:queueID/manifest", mainNode.Hostname))
	if err != nil {
		return nil, err
	}
	manifestResponse := new(http2.GetManifestResponse)

	err = json.NewDecoder(resp.Body).Decode(manifestResponse)
	if err != nil {
		return nil, err
	}

	return manifestResponse, nil
}

func (controller *Controller) syncFiles(ctx context.Context, mainNode *node) error {
	manifest, err := mainNode.getManifest(ctx)
	if err != nil {
		return err
	}

	for _, file := range manifest.Files {
		err = controller.syncFile(ctx, mainNode, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) syncFile(ctx context.Context, mainNode *node, fileID string) error {
	reader, err := mainNode.getReader(fileID)
	if err != nil {
		return err
	}

	writer, err := controller.getWriter(fileID)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	return nil
}

func (*Controller) getWriter(fileID string) (io.WriteCloser, error) {
	file, err := os.OpenFile(fmt.Sprintf("data1/%s", fileID), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	return file, err
}

func (mainNode *node) getReader(fileID string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/internal/queue/:queueID/file/%s", mainNode.Hostname, fileID))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return resp.Body, err
}

func (controller *Controller) deregister(ctx context.Context) error {
	_, err := controller.etcdClient.Delete(ctx, fmt.Sprintf("node/%s", controller.self.ID))
	return err
}

func (controller *Controller) Close() error {
	if controller.closed {
		return nil
	}

	controller.quit <- true
	return <-controller.err
}

func (controller *Controller) handleAsyncUpdates() {
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
