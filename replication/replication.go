package replication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	queueing "message-queueing"
	"net/http"
)

type repository struct {
	wrapped    queueing.Repository
	controller *Controller
}

func WrapRepository(repo queueing.Repository, controller *Controller) queueing.Repository {
	return &repository{
		wrapped:    repo,
		controller: controller,
	}
}

func (repo *repository) GetByID(ctx context.Context, id uuid.UUID) (*queueing.QueueMessage, error) {
	err := repo.controller.nodeService.DoOnAll(
		func(n *node, self bool, main bool) error {
			if self || main {
				return nil
			}
			return repo.syncGetByID(ctx, id, n)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get message by id due to upstream error: %w", err)
	}

	return repo.wrapped.GetByID(ctx, id)
}

func (repo *repository) syncGetByID(ctx context.Context, id uuid.UUID, n Node) error {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := repo.executeRequest(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("got illegal status code from upstream: %d", resp.StatusCode)
	}

	return nil
}

func (repo *repository) Create(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(
		func(n *node, self bool, main bool) error {
			if self || main {
				return nil
			}
			return repo.syncCreate(ctx, message, n)
		},
	)
	if err != nil {
		return fmt.Errorf("cannot create message due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Create call fails
	return repo.wrapped.Create(ctx, message)
}

func (repo *repository) syncCreate(ctx context.Context, message *queueing.QueueMessage, n Node) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/internal/queues/%s/messages", n.Host(), "abc"),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	resp, err := repo.executeRequest(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("got illegal status code from upstream: %d", resp.StatusCode)
	}

	return nil
}

func (repo *repository) Update(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(
		func(n *node, self bool, main bool) error {
			if self || main {
				return nil
			}
			return repo.syncUpdate(ctx, message, n)
		},
	)
	if err != nil {
		return fmt.Errorf("cannot update message due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Update call fails
	return repo.wrapped.Update(ctx, message)
}

func (repo *repository) syncUpdate(ctx context.Context, message *queueing.QueueMessage, n Node) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	id, err := uuid.FromBytes(message.MessageID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	resp, err := repo.executeRequest(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("got illegal status code from upstream: %d", resp.StatusCode)
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(
		func(n *node, self bool, main bool) error {
			if self || main {
				return nil
			}
			return repo.syncDelete(ctx, message, n)
		},
	)
	if err != nil {
		return fmt.Errorf("cannot get message by ID due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Delete call fails
	return repo.wrapped.Delete(ctx, message)
}

func (repo *repository) syncDelete(ctx context.Context, message *queueing.QueueMessage, n Node) error {
	id, err := uuid.FromBytes(message.MessageID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := repo.executeRequest(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("got illegal status code from upstream: %d", resp.StatusCode)
	}

	return nil
}

func (repo *repository) executeRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
}
