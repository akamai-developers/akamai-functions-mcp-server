package tools

import "log"

type AkamaiFunctionsTools struct {
	logger *log.Logger
}

func NewAkamaiFunctionsTools(logger *log.Logger) *AkamaiFunctionsTools {
	return &AkamaiFunctionsTools{
		logger: logger,
	}
}

type ToolResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func NewToolErrorResponse[T any](message string) ToolResponse[T] {
	return ToolResponse[T]{
		Success: false,
		Message: message,
	}
}

func NewToolSuccessResponse[T any](data T) ToolResponse[T] {
	return ToolResponse[T]{
		Success: true,
		Data:    data,
	}
}
