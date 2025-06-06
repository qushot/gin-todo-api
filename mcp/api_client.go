package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TodoAPIClient struct {
	baseURL string
	client  *http.Client
}

func NewTodoAPIClient(baseURL string) *TodoAPIClient {
	return &TodoAPIClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *TodoAPIClient) ListTodos(ctx context.Context) ([]Todo, error) {
	url := fmt.Sprintf("%s/todos", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return todos, nil
}

func (c *TodoAPIClient) GetTodo(ctx context.Context, id string) (*Todo, error) {
	url := fmt.Sprintf("%s/todos/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var todo Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &todo, nil
}

func (c *TodoAPIClient) CreateTodo(ctx context.Context, todo TodoRequestBody) (*Todo, error) {
	url := fmt.Sprintf("%s/todos", c.baseURL)

	jsonData, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var createdTodo Todo
	if err := json.NewDecoder(resp.Body).Decode(&createdTodo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdTodo, nil
}

func (c *TodoAPIClient) UpdateTodo(ctx context.Context, id string, todo TodoRequestBody) (*Todo, error) {
	url := fmt.Sprintf("%s/todos/%s", c.baseURL, id)

	jsonData, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var updatedTodo Todo
	if err := json.NewDecoder(resp.Body).Decode(&updatedTodo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedTodo, nil
}

func (c *TodoAPIClient) DeleteTodo(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/todos/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	return nil
}
