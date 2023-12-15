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
	nodes      []*node
	self       *node
	main       *node
	nodeChan   clientv3.WatchChan
	quit       chan bool
	err        chan error
	closed     bool

	dataDir string
}

func Open(ctx context.Context, etcdClient *clientv3.Client, dataDir string) (*Controller, error) {
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

	controller := &Controller{
		quit:       make(chan bool, 1),
		err:        make(chan error),
		nodes:      []*node{&self},
		self:       &self,
		main:       &self,
		etcdClient: etcdClient,
		closed:     false,
		dataDir:    dataDir,
	}

	err := controller.fetchNodes(ctx)
	if err != nil {
		return nil, err
	}

	if !controller.self.Main {
		err = controller.syncManifest(ctx)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("registering this node", "nodeID", fmt.Sprintf("node/%s", controller.self.ID))
	err = controller.register(ctx)
	if err != nil {
		return nil, err
	}

	controller.nodeChan = etcdClient.Watch(ctx, "node", clientv3.WithPrefix())

	go controller.handleAsyncUpdates()

	return controller, nil
}

func (controller *Controller) register(ctx context.Context) error {
	data, err := json.Marshal(controller.self)
	if err != nil {
		return err
	}

	_, err = controller.etcdClient.Put(ctx, fmt.Sprintf("node/%s", controller.self.ID), string(data))
	return err
}

func (controller *Controller) deregister(ctx context.Context) error {
	_, err := controller.etcdClient.Delete(ctx, fmt.Sprintf("node/%s", controller.self.ID))
	return err
}

func (controller *Controller) fetchNodes(ctx context.Context) error {
	resp, err := controller.etcdClient.Get(ctx, "node", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		n := new(node)
		err = json.Unmarshal(kv.Value, n)
		if err != nil {
			slog.ErrorContext(ctx, "could not parse node", "nodeID", kv.Key)
		}

		controller.nodes = append(controller.nodes, n)
		if n.Main {
			controller.main = n
		}
	}

	return nil
}

func (controller *Controller) syncManifest(ctx context.Context) error {
	manifest, err := controller.getManifest(ctx)
	if err != nil {
		return err
	}

	for _, blobID := range manifest.Blobs {
		err = controller.syncBlob(ctx, blobID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) getManifest(ctx context.Context) (*http2.GetManifestResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/internal/queues/abc/manifest", controller.main.Hostname))
	if err != nil {
		return nil, err
	}

	manifestResponse := new(http2.GetManifestResponse)
	err = json.NewDecoder(resp.Body).Decode(manifestResponse)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return manifestResponse, nil
}

func (controller *Controller) syncBlob(ctx context.Context, blobID string) error {
	reader, err := controller.getBlobReader(ctx, blobID)
	if err != nil {
		return err
	}

	writer, err := controller.getBlobWriter(ctx, blobID)
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

func (controller *Controller) getBlobReader(ctx context.Context, blobID string) (io.Reader, error) {
	resp, err := http.Get(
		fmt.Sprintf(
			"http://%s:8080/internal/queues/abc/blob/%s", controller.main.Hostname, blobID,
		),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return resp.Body, err
}

func (controller *Controller) getBlobWriter(ctx context.Context, blobID string) (io.WriteCloser, error) {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", controller.dataDir, blobID), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	return file, err
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

			n := new(node)
			err := json.Unmarshal(data.Kv.Value, n)
			if err != nil {
				slog.Error("could not unmarshal node", "nodeID", data.Kv.Key, "err", err)
				continue
			}

			slog.Info("received node update", "nodeID", data.Kv.Key)

			if data.IsCreate() {
				controller.addNode(n)
			} else if data.IsModify() {
				controller.updateNode(n)
			}
		case <-controller.quit:
			slog.Info("registering this node", "nodeID", fmt.Sprintf("node/%s", controller.self.ID))

			err := controller.deregister(context.TODO())
			controller.err <- err

			close(controller.quit)
			close(controller.err)
		}
	}
}

func (controller *Controller) addNode(newNode *node) {
	controller.nodes = append(controller.nodes, newNode)
}

func (controller *Controller) updateNode(updatedNode *node) {
	for i, n := range controller.nodes {
		if n.ID == updatedNode.ID {
			controller.nodes[i] = updatedNode
			return
		}
	}

	controller.addNode(updatedNode)
}
