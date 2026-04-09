package llm

// ChatCompletionRequest matches the BigModel "对话补全" request body (subset).
type ChatCompletionRequest struct {
	Model          string             `json:"model"`
	Messages       []ChatMessage      `json:"messages"`
	Temperature    float64            `json:"temperature,omitempty"`
	Stream         bool               `json:"stream"`
	ResponseFormat *ResponseFormat    `json:"response_format,omitempty"`
}

// ResponseFormat enables JSON mode when Type is "json_object".
type ResponseFormat struct {
	Type string `json:"type"`
}

// ChatMessage is one message in the conversation.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse is a minimal parse of the success JSON.
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
