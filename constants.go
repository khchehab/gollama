package gollama

type MessageRole string

type ThinkLevel string

type FormatType string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"

	ThinkHigh   ThinkLevel = "high"
	ThinkMedium ThinkLevel = "medium"
	ThinkLow    ThinkLevel = "low"

	FormatTypeJSON FormatType = "json"
)
