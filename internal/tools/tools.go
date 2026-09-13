package tools

import (
	"context"
	"encoding/json"
)

type Tool interface {
	Name() string
	Description() string
	InputSchema() json.RawMessage
	Execute(ctx context.Context, args json.RawMessage) (Result, error)
}

type Result struct {
	Content string
	IsError bool
}
