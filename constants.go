package gollama

type MessageRole string

type ChatMessageRole string

type ThinkLevel string

type FormatType string

type ToolType string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"

	ChatMessageRoleAssistant ChatMessageRole = "assistant"

	ThinkHigh   ThinkLevel = "high"
	ThinkMedium ThinkLevel = "medium"
	ThinkLow    ThinkLevel = "low"

	FormatTypeJSON FormatType = "json"

	ToolTypeFunction ToolType = "function"
)
