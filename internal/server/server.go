package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
	"github.com/incident-io/incidentio-mcp-golang/internal/handlers"
	"github.com/incident-io/incidentio-mcp-golang/pkg/mcp"
)

type Server struct {
	tools map[string]handlers.Handler
}

func New() *Server {
	return &Server{
		tools: make(map[string]handlers.Handler),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.registerTools()

	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var msg mcp.Message
			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF {
					return nil
				}
				continue
			}

			response, err := s.handleMessage(&msg)
			if err != nil {
				response = s.createErrorResponse(msg.ID, err)
			}

			if response != nil {
				if err := encoder.Encode(response); err != nil {
					// Log encoding errors but continue processing
					fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
				}
			}
		}
	}
}

func (s *Server) registerTools() {
	// Initialize incident.io client
	c, err := client.NewClient()
	if err != nil {
		// If client initialization fails, no tools are registered
		return
	}

	// Use the new registry to register all tools
	registry := handlers.NewToolRegistry()
	registry.RegisterAllTools(c)

	// Copy tools from registry to server
	s.tools = registry.GetTools()
}

func (s *Server) handleMessage(msg *mcp.Message) (*mcp.Message, error) {
	// Handle notifications (no ID means it's a notification)
	if msg.ID == nil {
		// Notifications don't require a response
		return nil, nil
	}

	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		// Return proper JSON-RPC error for unknown methods
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			},
		}, nil
	}
}

func (s *Server) handleInitialize(msg *mcp.Message) (*mcp.Message, error) {
	response := &mcp.Message{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "incidentio-mcp-server",
				"version": "0.1.0",
			},
		},
	}
	return response, nil
}

func (s *Server) handleToolsList(msg *mcp.Message) (*mcp.Message, error) {
	var toolsList []map[string]interface{}
	for _, tool := range s.tools {
		toolsList = append(toolsList, map[string]interface{}{
			"name":        tool.Name(),
			"description": tool.Description(),
			"inputSchema": tool.InputSchema(),
		})
	}

	response := &mcp.Message{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"tools": toolsList,
		},
	}
	return response, nil
}

func (s *Server) handleToolCall(msg *mcp.Message) (*mcp.Message, error) {
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing tool name")
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	args, _ := params["arguments"].(map[string]interface{})

	result, err := tool.Execute(args)
	if err != nil {
		return nil, err
	}

	response := &mcp.Message{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
	return response, nil
}

func (s *Server) createErrorResponse(id interface{}, err error) *mcp.Message {
	return &mcp.Message{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &mcp.Error{
			Code:    -32603,
			Message: err.Error(),
		},
	}
}
