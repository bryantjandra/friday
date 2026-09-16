package tools

import "github.com/bryantjandra/friday/internal/llm"

type Registry struct {
	ToolMap map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		ToolMap: make(map[string]Tool),
	}
}

func (r *Registry) RegisterTool(tool Tool) {
	r.ToolMap[tool.Name()] = tool
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, ok := r.ToolMap[name]
	if !ok {
		return nil, false
	}
	return tool, true
}

func (r *Registry) Definitions() []llm.ToolDef {
	/* we return an array of every single tool's name, description, and inputschema */
	var toolDefs []llm.ToolDef
	for _, tool := range r.ToolMap {
		toolDef := llm.ToolDef{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		}
		toolDefs = append(toolDefs, toolDef)
	}
	return toolDefs
}
