package llm

import "encoding/json"

/*
	Wire types: Go structs that match the JSON in the request and response. These types are used such that our Go code can use normal typed values
	and not raw JSON whenever working with the API.
*/

type Request struct {
	/* Marshalling: Turning a Go struct into JSON for requests */
	/*
		Struct Tags: By convention, a struct's fields should be capitalized (e.g. MaxTokens), but JSON values can be lowercase and of a different shape (e.g. max_tokens).
		Thus, we use these struct tags (backticks) to tell "encoding/json", that for example if it sees max_tokens in the JSON, then put it in the MaxTokens field, and vice versa.
	*/
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Tools     []ToolDef `json:"tools,omitempty"`
	Messages  []Message `json:"messages"`
	Thinking  *Thinking `json:"thinking,omitempty"`
}

type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

type Response struct {
	/* Unmarshalling: Turning JSON (from a response) into a Go struct */
	ID         string         `json:"id"`
	Model      string         `json:"model"`
	Role       string         `json:"role"`
	StopReason string         `json:"stop_reason"`
	Content    []ContentBlock `json:"content"`
	Usage      Usage          `json:"usage"`
}

type ContentBlock struct {
	Type      string          `json:"type,omitempty"`
	ID        string          `json:"id,omitempty"`
	Text      string          `json:"text,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   string          `json:"content,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Thinking struct {
	Type string `json:"type"`
}
