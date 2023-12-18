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

func (repo *repository) sendRequest(ctx context.Context, createRequest func(Node) (*http.Request, error)) DoOperation {
	return func(n *node, self bool, main bool) error {
		if self || main {
			return nil
		}
		request, err := createRequest(n)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("got illegal status code from upstream: %d", resp.StatusCode)
		}

		return nil
	}
}

func (repo *repository) GetByID(ctx context.Context, id uuid.UUID) (*queueing.QueueMessage, error) {
	err := repo.controller.nodeService.DoOnAll(repo.sendRequest(ctx, repo.buildGetByIDRequest(id)))
	if err != nil {
		return nil, fmt.Errorf("cannot get message by id due to upstream error: %w", err)
	}

	return repo.wrapped.GetByID(ctx, id)
}

func (repo *repository) buildGetByIDRequest(id uuid.UUID) func(Node) (*http.Request, error) {
	return func(n Node) (*http.Request, error) {
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
			nil,
		)
		if err != nil {
			return nil, err
		}

		return req, nil
	}
}

func (repo *repository) Create(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(repo.sendRequest(ctx, repo.buildCreateRequest(message)))
	if err != nil {
		return fmt.Errorf("cannot create message due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Create call fails
	return repo.wrapped.Create(ctx, message)
}

func (repo *repository) buildCreateRequest(message *queueing.QueueMessage) func(Node) (*http.Request, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	return func(n Node) (*http.Request, error) {
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://%s/internal/queues/%s/messages", n.Host(), "abc"),
			bytes.NewReader(data),
		)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return req, nil
	}
}

func (repo *repository) Update(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(repo.sendRequest(ctx, repo.createUpdateRequest(message)))
	if err != nil {
		return fmt.Errorf("cannot update message due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Update call fails
	return repo.wrapped.Update(ctx, message)
}

func (repo *repository) createUpdateRequest(message *queueing.QueueMessage) func(Node) (*http.Request, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	id, err := uuid.FromBytes(message.MessageID)
	if err != nil {
		return nil
	}

	return func(n Node) (*http.Request, error) {
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
			bytes.NewReader(data),
		)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

		return req, nil
	}
}

func (repo *repository) Delete(ctx context.Context, message *queueing.QueueMessage) error {
	err := repo.controller.nodeService.DoOnAll(repo.sendRequest(ctx, repo.createDeleteRequest(message)))
	if err != nil {
		return fmt.Errorf("cannot get message by ID due to upstream error: %w", err)
	}

	// TODO: implement rollback when upstream Delete call fails
	return repo.wrapped.Delete(ctx, message)
}

func (repo *repository) createDeleteRequest(message *queueing.QueueMessage) func(Node) (*http.Request, error) {
	id, err := uuid.FromBytes(message.MessageID)
	if err != nil {
		return nil
	}

	return func(n Node) (*http.Request, error) {
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("http://%s/internal/queues/%s/messages/%s", n.Host(), "abc", id.String()),
			nil,
		)
		if err != nil {
			return nil, err
		}

		return req, nil
	}
}
