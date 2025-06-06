package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type TodoRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type Todo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"ToDo API 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	apiClient := NewTodoAPIClient("http://localhost:8080/api/v1")

	registerListTodosTool(s, apiClient)
	registerReadTodoTool(s, apiClient)
	registerCreateTodoTool(s, apiClient)
	registerUpdateTodoTool(s, apiClient)
	registerDeleteTodoTool(s, apiClient)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func registerCreateTodoTool(s *server.MCPServer, apiClient *TodoAPIClient) {
	createTodoTool := mcp.NewTool("create_todo",
		mcp.WithDescription("Create a new todo"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the todo item"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content of the todo item"),
		),
		mcp.WithBoolean("done",
			mcp.Required(),
			mcp.Description("Is the todo item done?"),
		),
	)
	createTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		done, err := request.RequireBool("done")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo := TodoRequestBody{
			Title:   title,
			Content: content,
			Done:    done,
		}

		createdTodo, err := apiClient.CreateTodo(ctx, todo)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Created todo: %s - %s (done: %v)", createdTodo.Title, createdTodo.Content, createdTodo.Done)), nil
	}
	s.AddTool(createTodoTool, createTodoHandler)
}

func registerListTodosTool(s *server.MCPServer, apiClient *TodoAPIClient) {
	listTodosTool := mcp.NewTool("list_todos",
		mcp.WithDescription("List all todos"),
	)
	listTodosHandler := func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		todos, err := apiClient.ListTodos(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		resultText := "Todos:\n"
		for _, todo := range todos {
			resultText += fmt.Sprintf("- [%s] %s: %s (done: %v)\n", todo.ID, todo.Title, todo.Content, todo.Done)
		}

		return mcp.NewToolResultText(resultText), nil
	}
	s.AddTool(listTodosTool, listTodosHandler)
}

func registerReadTodoTool(s *server.MCPServer, apiClient *TodoAPIClient) {
	readTodoTool := mcp.NewTool("read_todo",
		mcp.WithDescription("Read a todo by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the todo item to read"),
		),
	)
	readTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo, err := apiClient.GetTodo(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		resultText := fmt.Sprintf("Todo: [%s] %s: %s (done: %v)", todo.ID, todo.Title, todo.Content, todo.Done)
		return mcp.NewToolResultText(resultText), nil
	}
	s.AddTool(readTodoTool, readTodoHandler)
}

func registerUpdateTodoTool(s *server.MCPServer, apiClient *TodoAPIClient) {
	updateTodoTool := mcp.NewTool("update_todo",
		mcp.WithDescription("Update a todo by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the todo item to update"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("New title of the todo item"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("New content of the todo item"),
		),
		mcp.WithBoolean("done",
			mcp.Required(),
			mcp.Description("Is the todo item done?"),
		),
	)
	updateTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		done, err := request.RequireBool("done")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo := TodoRequestBody{
			Title:   title,
			Content: content,
			Done:    done,
		}

		updatedTodo, err := apiClient.UpdateTodo(ctx, id, todo)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Updated todo with ID %s: %s - %s (done: %v)", id, updatedTodo.Title, updatedTodo.Content, updatedTodo.Done)), nil
	}
	s.AddTool(updateTodoTool, updateTodoHandler)
}

func registerDeleteTodoTool(s *server.MCPServer, apiClient *TodoAPIClient) {
	deleteTodoTool := mcp.NewTool("delete_todo",
		mcp.WithDescription("Delete a todo by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the todo item to delete"),
		),
	)
	deleteTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := apiClient.DeleteTodo(ctx, id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Deleted todo with ID %s", id)), nil
	}
	s.AddTool(deleteTodoTool, deleteTodoHandler)
}
