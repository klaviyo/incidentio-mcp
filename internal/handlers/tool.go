package handlers

// Handler represents an MCP tool handler
type Handler interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}
	Execute(args map[string]interface{}) (string, error)
}
