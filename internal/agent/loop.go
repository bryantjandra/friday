package agent

import (
	"context"
	"errors"
	"strings"

	"github.com/bryantjandra/friday/internal/llm"
	"github.com/bryantjandra/friday/internal/tools"
)

type Agent struct {
	client   *llm.Client
	registry *tools.Registry
}

func NewAgent(client *llm.Client, registry *tools.Registry) *Agent {
	return &Agent{
		client:   client,
		registry: registry,
	}
}

/*

- When `Run` is called, it first builds the **system prompt** (carrying today's date and timezone, since the model has no clock).
- It then constructs the **message history** — a `messages` array that starts with a single entry: the user's initial prompt, tagged with `role: "user"`.
- Next it builds the **request**, which bundles: the max tokens, the system prompt, the message history, the definitions of all available tools, and thinking set to `"disabled"`.

Then it enters a loop (capped at 10 iterations to guarantee termination). Each iteration:
1. **Call the model** with the current request.
2. **Append the model's response to the history.** Specifically, the response's entire `content` array (which may hold multiple blocks — text, thinking, and/or tool_use) is appended as a new message tagged `role: "assistant"`.
This is essential because the model is **stateless** — it remembers nothing between calls, so the full conversation must be resupplied every time. Appending the assistant's turn (including any tool_use blocks and their ids) is what gives the model context across iterations.
3. **Branch on the stop reason:**
   - **If `stop_reason == "tool_use"`**, the model is asking to run one or more tools. Walk the `content` array; for each block of type `tool_use`, look the tool up in the registry by name and execute it with the block's `input`. Collect each outcome into a `toolResults` array, where every entry carries the `tool_use_id` (matching the call it answers) and a `content` string — either a success message or an error describing what went wrong (e.g. a bad date).
   Once all tool_use blocks are handled, append `toolResults` to the history as a new message tagged `role: "user"`. *(Note: tool results are tagged `user` not because a human produced them, but because the API treats "the model" as the assistant and everything else — including tool output — as the user side.)* Then `continue` to the next iteration, so the model can see the results and decide its next move.
   - **If `stop_reason` is anything else** (normally `end_turn` — the model has finished), collect the text from the `content` array's text blocks, join them, and return that as the final answer.

*/

func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {

	systemPrompt := BuildSystemPrompt()

	messageHistory := []llm.Message{
		{
			Role:    "user",
			Content: []llm.ContentBlock{{Type: "text", Text: prompt}},
		},
	}
	req := llm.Request{
		MaxTokens: 1024,
		Messages:  messageHistory,
		System:    systemPrompt,
		Tools:     a.registry.Definitions(),
		Thinking:  &llm.Thinking{Type: "disabled"},
	}

	for i := 0; i < 10; i++ {

		resp, err := a.client.CreateMessage(ctx, req)

		if err != nil {
			return "", err
		}

		req.Messages = append(req.Messages, llm.Message{Role: "assistant", Content: resp.Content})

		if resp.StopReason == "tool_use" {
			toolResults := []llm.ContentBlock{}
			for _, block := range resp.Content {
				if block.Type == "tool_use" {
					tool, ok := a.registry.GetTool(block.Name)
					if !ok {
						return "", errors.New("invalid tool was called")
					}
					res, err := tool.Execute(ctx, block.Input)
					if err != nil {
						return "", err
					}
					toolResults = append(toolResults, llm.ContentBlock{Type: "tool_result",
						ToolUseID: block.ID,
						Content:   res.Content,
						IsError:   res.IsError})

				}
			}
			req.Messages = append(req.Messages, llm.Message{Role: "user", Content: toolResults})
			continue
		}
		res := []string{}
		for _, block := range resp.Content {
			res = append(res, block.Text)
		}

		resSingle := strings.Join(res, "")
		return resSingle, nil
	}

	return "", errors.New("agent exceeded maximum amount of iterations before completing")
}
